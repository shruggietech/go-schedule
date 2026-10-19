package clientsecret

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const readHelperArgument = "--internal-native-credential-read"
const readResponseLimit = 64 * 1024

var errReadUnavailable = errors.New("native credential read is unavailable; check the credential service or repair this connection")

type contextBackend interface {
	GetContext(context.Context, string, string) (string, error)
}
type readRequest struct {
	DaemonID     string `json:"daemon_id"`
	CredentialID string `json:"credential_id"`
}
type readResponse struct {
	Token       string `json:"token,omitempty"`
	Unavailable bool   `json:"unavailable,omitempty"`
}

// LoadContext reads within the caller's deadline, joining its owned helper on cancellation.
// Injected backends must provide GetContext; a blocking legacy Get is never called.
func (s *Store) LoadContext(ctx context.Context, daemonID, credentialID string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	backend, ok := s.backend.(contextBackend)
	if !ok {
		return "", errors.New("credential backend does not support cancellation-aware reads")
	}
	return backend.GetContext(ctx, servicePrefix+daemonID, credentialID)
}
func (nativeBackend) GetContext(ctx context.Context, service, user string) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", errReadUnavailable
	}
	return runReadCommand(ctx, executable, []string{readHelperArgument}, readRequest{DaemonID: strings.TrimPrefix(service, servicePrefix), CredentialID: user})
}

// RunReadHelper handles only the private native-read subprocess protocol.
// CLI hosts call it before their normal entry point. Secrets travel only over redirected pipes.
func RunReadHelper(args []string, input io.Reader, output io.Writer) (bool, int) {
	if len(args) != 1 || args[0] != readHelperArgument {
		return false, 0
	}
	return true, serveReadHelper(input, output, nativeBackend{})
}
func serveReadHelper(input io.Reader, output io.Writer, backend Backend) int {
	var request readRequest
	decoder := json.NewDecoder(io.LimitReader(input, 16*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return 1
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF || strings.TrimSpace(request.DaemonID) == "" || strings.TrimSpace(request.CredentialID) == "" {
		return 1
	}
	token, err := backend.Get(servicePrefix+request.DaemonID, request.CredentialID)
	response := readResponse{Token: token}
	if err != nil {
		response = readResponse{Unavailable: true}
	}
	if json.NewEncoder(output).Encode(response) != nil {
		return 1
	}
	return 0
}

type boundedResponse struct{ data bytes.Buffer }

func (b *boundedResponse) Write(data []byte) (int, error) {
	if b.Len()+len(data) > readResponseLimit {
		return 0, errReadUnavailable
	}
	return b.data.Write(data)
}
func runReadCommand(ctx context.Context, executable string, args []string, request readRequest) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	input, err := json.Marshal(request)
	if err != nil {
		return "", errReadUnavailable
	}
	command := exec.CommandContext(ctx, executable, args...)
	command.Stdin = bytes.NewReader(input)
	var output boundedResponse
	command.Stdout = &output
	command.Stderr = io.Discard
	command.WaitDelay = 250 * time.Millisecond
	cleanup := configureReadCommand(command)
	runErr := command.Run()
	cleanupErr := cleanup()
	if err := ctx.Err(); err != nil {
		return "", errors.Join(err, cleanupErr)
	}
	if runErr != nil || cleanupErr != nil {
		return "", errReadUnavailable
	}
	var response readResponse
	if json.Unmarshal(output.Bytes(), &response) != nil || response.Unavailable || response.Token == "" {
		return "", errReadUnavailable
	}
	return response.Token, nil
}

func (b *boundedResponse) Len() int      { return b.data.Len() }
func (b *boundedResponse) Bytes() []byte { return b.data.Bytes() }
