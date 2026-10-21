package clientsecret

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
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

// RunReadHelper handles only the private, peer-authenticated native-read channel.
// Standard input and output are never credential protocol endpoints.
func RunReadHelper(args []string) (bool, int) {
	if len(args) == 0 || args[0] != readHelperArgument {
		return false, 0
	}
	if len(args) != 2 {
		return true, 1
	}
	return true, runAuthenticatedHelper(args[1], nativeBackend{}, 15*time.Second)
}

// An independent process deadline also ends a stalled read if its parent dies.
func runAuthenticatedHelper(endpoint string, backend Backend, limit time.Duration) int {
	timer := time.AfterFunc(limit, expireReadHelper)
	defer timer.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	return serveAuthenticatedHelper(ctx, endpoint, backend)
}
func serveAuthenticatedHelper(ctx context.Context, endpoint string, backend Backend) int {
	connection, err := dialReadChannel(ctx, endpoint)
	if err != nil {
		return 1
	}
	return serveAuthenticatedConnection(ctx, connection, backend)
}

// serveAuthenticatedConnection owns an established channel and authenticates its
// native peer before reading requests or accessing the credential backend.
func serveAuthenticatedConnection(ctx context.Context, connection net.Conn, backend Backend) int {
	defer func() { _ = connection.Close() }()
	stop := context.AfterFunc(ctx, func() { _ = connection.Close() })
	defer stop()
	peer, err := readPeerPID(connection, false)
	if err != nil || authenticateReadParent(peer) != nil {
		return 1
	}
	return serveReadHelper(connection, connection, backend)
}
func authenticateReadParent(peer int) error {
	if peer <= 1 || peer != os.Getppid() {
		return errReadUnavailable
	}
	path, err := readProcessExecutable(peer)
	if err != nil {
		return errReadUnavailable
	}
	self, err := os.Executable()
	if err != nil {
		return errReadUnavailable
	}
	return sameReadExecutable(self, path)
}
func sameReadExecutable(self, parent string) error {
	own, err := os.Stat(self)
	if err != nil {
		return errReadUnavailable
	}
	other, err := os.Stat(parent)
	if err != nil || !os.SameFile(own, other) {
		return errReadUnavailable
	}
	return nil
}
func serveReadHelper(input io.Reader, output io.Writer, backend Backend) int {
	data, err := readFrame(input, 16*1024)
	if err != nil {
		return 1
	}
	var request readRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&request) != nil {
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
	encoded, err := json.Marshal(response)
	if err != nil || len(encoded) > readResponseLimit {
		return 1
	}
	if writeFrame(output, encoded) != nil {
		return 1
	}
	return 0
}
func readFrame(input io.Reader, limit int) ([]byte, error) {
	var size uint32
	if binary.Read(input, binary.BigEndian, &size) != nil || size == 0 || uint64(size) > uint64(limit) {
		return nil, errReadUnavailable
	}
	data := make([]byte, int(size))
	if _, err := io.ReadFull(input, data); err != nil {
		return nil, errReadUnavailable
	}
	return data, nil
}
func writeFrame(output io.Writer, data []byte) error {
	if err := binary.Write(output, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := io.Copy(output, bytes.NewReader(data))
	return err
}
func acceptReadHelper(ctx context.Context, listener net.Listener, pid int) (net.Conn, error) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, errReadUnavailable
		}
		peer, err := readPeerPID(connection, true)
		if err == nil && peer == pid {
			return connection, nil
		}
		_ = connection.Close()
	}
}
func runReadCommand(ctx context.Context, executable string, args []string, request readRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	listener, endpoint, closeChannel, err := listenReadChannel()
	if err != nil {
		return "", errReadUnavailable
	}
	defer closeChannel()
	stopAccept := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stopAccept()
	command := exec.CommandContext(ctx, executable, append(append([]string{}, args...), endpoint)...)
	command.Stdin = nil
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	command.WaitDelay = 250 * time.Millisecond
	cleanup := configureReadCommand(command)
	if command.Start() != nil {
		return "", errReadUnavailable
	}
	finished := make(chan error, 1)
	go func() { finished <- command.Wait(); _ = listener.Close() }()
	connection, exchangeErr := acceptReadHelper(ctx, listener, command.Process.Pid)
	var response readResponse
	if exchangeErr == nil {
		stopConnection := context.AfterFunc(ctx, func() { _ = connection.Close() })
		encoded, _ := json.Marshal(request)
		exchangeErr = writeFrame(connection, encoded)
		if exchangeErr == nil {
			var data []byte
			data, exchangeErr = readFrame(connection, readResponseLimit)
			if exchangeErr == nil {
				exchangeErr = json.Unmarshal(data, &response)
			}
		}
		_ = connection.Close()
		stopConnection()
	}
	requestErr := ctx.Err()
	if exchangeErr != nil {
		cancel()
	}
	runErr := <-finished
	cleanupErr := cleanup()
	if requestErr == nil && exchangeErr == nil {
		requestErr = ctx.Err()
	}
	if requestErr != nil {
		return "", errors.Join(requestErr, cleanupErr)
	}
	if exchangeErr != nil || runErr != nil || cleanupErr != nil || response.Unavailable || response.Token == "" {
		return "", errReadUnavailable
	}
	return response.Token, nil
}
