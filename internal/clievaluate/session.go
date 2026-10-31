package clievaluate

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/cli"
	"github.com/shruggietech/go-schedule/internal/clientprofile"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/incident"
	"github.com/shruggietech/go-schedule/internal/remoteenroll"
	"github.com/shruggietech/go-schedule/internal/testsmtp"
)

type session struct {
	smtp                 *testsmtp.Server
	smtpRotation         domain.SMTPAuth
	smtpBeforeCount      int
	priorRuns            map[string]bool
	w                    *clifixture.Workspace
	api                  *client.Client
	binaries             map[string]string
	vars                 map[string]string
	redactor             *Redactor
	budget               Budget
	recorder             *clifixture.Recorder
	profiles             *clientprofile.Store
	daemon               *clifixture.Process
	incidents            *incident.Channel
	incidentError        error
	daemonPipes          *clifixture.Pipes
	daemonOut, daemonErr *memoryCapture
	daemonDone           chan error
	webhook              *http.Server
	webhookDone          chan error
	webhookMu            sync.Mutex
	webhookHits          int
	webhookCode          int
}

func newSession(ctx context.Context, binaries map[string]string, budget Budget, redactor *Redactor, incidents bool, onOwned func(*clifixture.Workspace) error) (s *session, resultErr error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		return nil, err
	}
	s = &session{w: w, binaries: binaries, budget: budget, redactor: redactor, recorder: &clifixture.Recorder{}, vars: map[string]string{}, webhookCode: 200}
	defer func() {
		if resultErr != nil {
			resultErr = errors.Join(resultErr, s.close())
		}
	}()
	if onOwned != nil {
		if err := onOwned(w); err != nil {
			return s, err
		}
	}
	s.vars["webhook_secret"] = uuid.NewString()
	s.vars["webhook_rotated"] = uuid.NewString()
	s.redactor.Register(s.vars["webhook_secret"], s.vars["webhook_rotated"])
	s.profiles = clientprofile.NewStore(filepath.Join(w.Runtime, "client", "profiles.json"))
	s.api = client.New(w.IPCPath)
	s.vars["ipc"] = w.IPCPath
	s.vars["missing_ipc"] = w.IPCPath + "-absent"
	s.vars["runtime"] = w.Runtime
	s.vars["config"] = w.ConfigPath
	s.vars["worker"] = binaries["runner"]
	s.vars["future"] = time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	s.vars["watch_root"] = filepath.Join(w.Runtime, "watchers")
	s.vars["stdin_file"] = filepath.Join(w.Runtime, "tasks", "stdin.txt")
	s.vars["cron_file"] = filepath.Join(w.Runtime, "tasks", "fixture.cron")
	s.vars["bad_json"] = filepath.Join(w.Runtime, "tasks", "bad.json")
	s.vars["bad_certificate"] = filepath.Join(w.Runtime, "certificates", "bad.pem")
	s.vars["policy_file"] = filepath.Join(w.Runtime, "tasks", "policy.json")
	s.vars["audit_file"] = filepath.Join(w.Runtime, "exports", "audit.jsonl")
	s.vars["bundle_file"] = filepath.Join(w.Runtime, "exports", "bundle.json")
	s.vars["plan_file"] = filepath.Join(w.Runtime, "exports", "plan.json")
	for key, value := range map[string]string{"stdin_file": "fixture stdin 世界\n", "cron_file": "0 0 1 1 * echo fixture\n", "bad_json": "{truncated", "bad_certificate": "not a certificate"} {
		if err := s.writeOwned(s.vars[key], []byte(value)); err != nil {
			return s, err
		}
	}
	certificate, keyPath, certificatePEM, err := makeCertificate(w.Runtime)
	if err != nil {
		return s, err
	}
	s.vars["certificate"] = certificate
	s.vars["certificate_pem"] = certificatePEM
	address, err := availableAddress()
	if err != nil {
		return s, err
	}
	s.vars["remote"] = "https://" + address
	w.Config.Remote.Enabled = true
	w.Config.Remote.BindAddress = address
	w.Config.Remote.CertificateFile = certificate
	w.Config.Remote.PrivateKeyFile = keyPath
	w.Config.LogFormat = "json"
	w.Config.DefaultTimezone = "UTC"
	configData, err := json.MarshalIndent(w.Config, "", "  ")
	if err != nil {
		return s, err
	}
	if err := s.writeOwned(w.ConfigPath, configData); err != nil {
		return s, err
	}
	s.daemonOut = newMemoryCapture(budget.ChannelBytes, nil)
	s.daemonErr = newMemoryCapture(budget.ChannelBytes, nil)
	environment := incident.Environment(quietEnvironment(), map[string]string{"GOSCHEDULE_QA_INCIDENT_DIR": ""})
	if incidents {
		s.incidents, err = incident.NewChannel(w.Evidence)
		if err == nil {
			err = w.Contains(s.incidents.Dir)
		}
		if err != nil {
			s.incidentError = err
			s.incidents = nil
		} else {
			environment = incident.Environment(environment, map[string]string{"GOSCHEDULE_QA_INCIDENT_DIR": s.incidents.Dir})
		}
	}
	s.daemon, s.daemonPipes, err = clifixture.StartPiped(clifixture.Command{Executable: binaries["daemon"], Arguments: []string{"--config", w.ConfigPath}, Directory: w.Runtime, Environment: environment})
	if err != nil {
		return s, err
	}
	if err := s.daemonPipes.Stdin.Close(); err != nil {
		return s, err
	}
	s.daemonDone = make(chan error, 2)
	go func() { _, err := io.Copy(s.daemonOut, s.daemonPipes.Stdout); s.daemonDone <- err }()
	go func() { _, err := io.Copy(s.daemonErr, s.daemonPipes.Stderr); s.daemonDone <- err }()
	if err := waitUntil(ctx, func() bool { _, err := s.api.Health(ctx); return err == nil }); err != nil {
		return s, fmt.Errorf("owned daemon readiness: %w", err)
	}
	if err := w.Bind(ctx, s.api, s.daemon.PID); err != nil {
		return s, err
	}
	if onOwned != nil {
		if err := onOwned(w); err != nil {
			return s, err
		}
	}
	s.vars["daemon_id"] = w.DaemonID
	currentUser, err := user.Current()
	if err != nil {
		return s, err
	}
	s.vars["current_user"] = currentUser.Username
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return s, err
	}
	s.vars["webhook"] = "http://" + listener.Addr().String() + "/fixture"
	s.webhook = &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		s.webhookMu.Lock()
		s.webhookHits++
		code := s.webhookCode
		s.webhookMu.Unlock()
		w.WriteHeader(code)
	}), ReadHeaderTimeout: 2 * time.Second}
	s.webhookDone = make(chan error, 1)
	go func() {
		err := s.webhook.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		s.webhookDone <- err
	}()
	if err := s.startSMTP(); err != nil {
		return s, err
	}
	if err := s.seed(ctx); err != nil {
		return s, err
	}
	return s, nil
}

func quietEnvironment() []string {
	return append(os.Environ(), "GOMAXPROCS=2", "GOFLAGS=-p=2", "GOMEMLIMIT=512MiB", "GIT_TERMINAL_PROMPT=0", "GCM_INTERACTIVE=Never")
}
func availableAddress() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	address := listener.Addr().String()
	return address, listener.Close()
}
func waitUntil(ctx context.Context, ready func() bool) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if ready() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func makeCertificate(root string) (string, string, string, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 120))
	if err != nil {
		return "", "", "", err
	}
	now := time.Now()
	template := &x509.Certificate{SerialNumber: serial, Subject: pkix.Name{CommonName: "owned CLI QA loopback"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(2 * time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return "", "", "", err
	}
	private, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", "", err
	}
	certificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: private})
	certificatePath, keyPath := filepath.Join(root, "certificates", "server.pem"), filepath.Join(root, "certificates", "server.key")
	if err := os.WriteFile(certificatePath, certificate, 0600); err != nil {
		return "", "", "", err
	}
	if err := os.WriteFile(keyPath, keyBytes, 0600); err != nil {
		return "", "", "", err
	}
	return certificatePath, keyPath, string(certificate), nil
}

func (s *session) writeOwned(path string, data []byte) error {
	if err := s.w.Contains(path); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (s *session) seed(ctx context.Context) error {
	group, err := s.api.CreateGroup(ctx, server.GroupCreateRequest{Name: "QA baseline"})
	if err != nil {
		return err
	}
	s.vars["base_group"] = group.ID
	for _, definition := range []struct{ variable, mode string }{{"base_task", "success"}, {"second_task", "success"}, {"failure_task", "failure"}, {"output_task", "output"}, {"active_task", "wait"}} {
		created, err := s.api.CreateTask(ctx, server.TaskCreateRequest{Name: "QA " + definition.variable, GroupID: group.ID, Command: s.binaries["runner"], Args: []string{"fixture-child", "--mode=" + definition.mode}, WorkingDir: s.w.Runtime, Env: map[string]string{"QA_VALUE": "owned 世界"}, Timezone: "UTC", Schedule: "every 24 hours"})
		if err != nil {
			return err
		}
		s.vars[definition.variable] = created.Task.ID
	}
	missing, err := s.api.CreateTask(ctx, server.TaskCreateRequest{Name: "QA missing executable", Command: filepath.Join(s.w.Runtime, "tasks", "missing-executable"), Timezone: "UTC"})
	if err != nil {
		return err
	}
	s.vars["missing_task"] = missing.Task.ID
	second, err := s.api.CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: "QA baseline webhook", Endpoint: s.vars["webhook"]})
	if err != nil {
		return err
	}
	s.vars["second_channel"] = second.ID
	secret, err := s.api.CreatePairing(ctx, server.PairingCreateRequest{DisplayName: "QA baseline remote", Kind: domain.ActorKindCLI, Capability: domain.CapabilityManage})
	if err != nil {
		return err
	}
	s.redactor.Register(secret.Phrase)
	credential, manifest, err := remoteenroll.Exchange(ctx, remoteenroll.Draft{Address: s.vars["remote"], DaemonID: secret.DaemonID, PairingID: secret.ID, Phrase: secret.Phrase, CertificatePEM: s.vars["certificate_pem"], DisplayName: secret.DisplayName, Kind: secret.Kind, Capability: secret.Capability})
	if err != nil {
		return err
	}
	s.redactor.Register(credential.Token)
	s.vars["remote_credential"] = credential.ID
	s.vars["remote_token"] = credential.Token
	if err := clientsecret.NewWithBackend(s.recorder).Save(manifest.InstallationID, credential.ID, credential.Token); err != nil {
		return err
	}
	profileID, err := clientprofile.NewID()
	if err != nil {
		return err
	}
	profile, err := s.profiles.Add(clientprofile.Profile{ID: profileID, ClientKind: "cli", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Label: "QA baseline", Endpoint: s.vars["remote"], DaemonID: manifest.InstallationID, CredentialID: credential.ID, CertificatePEM: s.vars["certificate_pem"], Capability: string(domain.CapabilityManage)})
	if err != nil {
		return err
	}
	s.vars["remote_profile"] = profile.ID
	third, err := s.api.CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: "QA alternate webhook", Endpoint: s.vars["webhook"]})
	if err != nil {
		return err
	}
	s.vars["third_channel"] = third.ID
	policy := server.NotificationAssignmentsRequest{Assignments: []server.NotificationAssignmentInput{{ChannelID: second.ID, OnSuccess: true}, {ChannelID: third.ID, OnFailure: true, FailureThreshold: 2, OnFailureToStart: true, OnRecovery: true}}}
	policyData, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	s.vars["policy_json"] = string(policyData)
	if err := s.writeOwned(s.vars["policy_file"], policyData); err != nil {
		return err
	}
	port, err := availableAddress()
	if err != nil {
		return err
	}
	_, number, err := net.SplitHostPort(port)
	if err != nil {
		return err
	}
	s.vars["mcp_port"] = number
	return nil
}

func (s *session) dependencies(in io.Reader, out, diagnostic io.Writer) cli.Dependencies {
	return cli.Dependencies{Client: s.api, Profiles: s.profiles, Secrets: clientsecret.NewWithBackend(s.recorder), PreferencePath: filepath.Join(s.w.Runtime, "client", "preferences.json"), Now: time.Now, ServiceControl: s.recorder.ServiceControl, DaemonPath: func() (string, error) { return s.binaries["daemon"], nil }, GUIPath: func() (string, error) { return s.binaries["runner"], nil }, LaunchGUI: s.recorder.LaunchGUI, In: in, Out: out, Err: diagnostic}
}

func (s *session) expand(value string) (string, error) {
	for {
		start := strings.Index(value, "${")
		if start < 0 {
			return value, nil
		}
		end := strings.IndexByte(value[start:], '}')
		if end < 0 {
			return "", errors.New("invalid fixture variable")
		}
		key := value[start+2 : start+end]
		replacement, ok := s.vars[key]
		if !ok {
			return "", fmt.Errorf("fixture prerequisite %s is unavailable", key)
		}
		value = value[:start] + replacement + value[start+end+1:]
	}
}
func (s *session) arguments(args []string) ([]string, error) {
	result := make([]string, len(args))
	for i, arg := range args {
		value, err := s.expand(arg)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

func (s *session) close() error {
	var result error
	if s.webhook != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		result = errors.Join(result, s.webhook.Shutdown(ctx))
		cancel()
		if s.webhookDone != nil {
			result = errors.Join(result, <-s.webhookDone)
		}
		s.webhook = nil
	}
	if s.daemon != nil {
		result = errors.Join(result, s.daemon.Stop())
		if result != nil && s.daemonPipes != nil {
			result = errors.Join(result, s.daemonPipes.Close())
		}
		if s.daemonDone != nil {
			for i := 0; i < 2; i++ {
				result = errors.Join(result, <-s.daemonDone)
			}
			s.daemonDone = nil
		}
		if s.daemonPipes != nil {
			result = errors.Join(result, s.daemonPipes.Close())
		}
		s.daemon = nil
	}
	if s.smtp != nil {
		result = errors.Join(result, s.smtp.Close())
		s.smtp = nil
	}
	if result == nil && s.w != nil {
		result = s.w.CleanupRuntime()
	}
	return result
}

func (s *session) snapshot(ctx context.Context) (map[string]any, error) {
	result := map[string]any{}
	calls := []struct {
		name string
		load func() (any, error)
	}{{"tasks", func() (any, error) { return s.api.ListTasks(ctx, "", "") }}, {"groups", func() (any, error) { return s.api.ListGroups(ctx) }}, {"chains", func() (any, error) { return s.api.ListChains(ctx) }}, {"triggers", func() (any, error) { return s.api.ListTriggers(ctx) }}, {"sets", func() (any, error) { return s.api.ListTriggerSets(ctx) }}, {"watchers", func() (any, error) { return s.api.ListFilesystemWatchers(ctx) }}, {"channels", func() (any, error) { return s.api.ListNotificationChannels(ctx) }}, {"actors", func() (any, error) { return s.api.ListActors(ctx) }}, {"credentials", func() (any, error) { return s.api.ListCredentials(ctx) }}, {"pairings", func() (any, error) { return s.api.ListPairings(ctx) }}}
	for _, call := range calls {
		value, err := call.load()
		if err != nil {
			return result, fmt.Errorf("snapshot %s: %w", call.name, err)
		}
		data, err := json.Marshal(value)
		if err != nil {
			return result, err
		}
		var items []map[string]any
		if err := json.Unmarshal(data, &items); err != nil {
			return result, err
		}
		indexed := map[string]any{}
		for _, item := range items {
			id, ok := item["id"].(string)
			if !ok {
				return result, fmt.Errorf("%s observation has no ID", call.name)
			}
			indexed[id] = item
		}
		result[call.name] = indexed
	}
	profiles, err := s.profiles.Load()
	if err != nil {
		return result, err
	}
	data, err := json.Marshal(profiles)
	if err != nil {
		return result, err
	}
	var profileValue any
	if err := json.Unmarshal(data, &profileValue); err != nil {
		return result, err
	}
	result["profiles"] = profileValue
	intents := []clifixture.Intent{}
	for _, intent := range s.recorder.Intents() {
		if strings.HasPrefix(intent.Operation, "service.") || intent.Operation == "gui.launch" {
			intents = append(intents, intent)
		}
	}
	result["host_intents"] = intents
	prefData, prefErr := os.ReadFile(filepath.Join(s.w.Runtime, "client", "preferences.json"))
	if prefErr != nil && !errors.Is(prefErr, os.ErrNotExist) {
		return result, prefErr
	}
	result["preferences"] = string(prefData)
	return result, nil
}

func fieldAt(value any, path string) (any, bool) {
	for _, part := range strings.Split(path, ".") {
		switch node := value.(type) {
		case map[string]any:
			next, ok := node[part]
			if !ok {
				return nil, false
			}
			value = next
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(node) {
				return nil, false
			}
			value = node[index]
		default:
			return nil, false
		}
	}
	return value, true
}
func printable(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	if value == nil {
		return "<nil>"
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "<invalid observation>"
	}
	return string(data)
}
