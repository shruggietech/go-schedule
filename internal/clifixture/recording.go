package clifixture

import (
	"context"
	"errors"
	"sync"
)

// Intent records an OS-facing operation without retaining secret values.
type Intent struct {
	Operation, Executable string
	Arguments             []string
	Evidence              EvidenceKind
}

// Recorder implements service, launch and credential seams without host effects.
type Recorder struct {
	mu      sync.Mutex
	intents []Intent
	values  map[string]string
	Failure error
}

// Intents returns a copy of the observed operation intents.
func (r *Recorder) Intents() []Intent {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Intent, len(r.intents))
	copy(out, r.intents)
	for i := range out {
		out[i].Arguments = append([]string(nil), out[i].Arguments...)
	}
	return out
}
func (r *Recorder) record(intent Intent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	intent.Evidence = CommandHarness
	r.intents = append(r.intents, intent)
	return r.Failure
}

// ServiceControl records an actual handler's requested service action.
func (r *Recorder) ServiceControl(action, executable string, args []string) (string, error) {
	if err := r.record(Intent{Operation: "service." + action, Executable: executable, Arguments: append([]string(nil), args...)}); err != nil {
		return "", err
	}
	return "recorded " + action, nil
}

// LaunchGUI records the intended launch without creating a desktop process.
func (r *Recorder) LaunchGUI(path string) error {
	return r.record(Intent{Operation: "gui.launch", Executable: path})
}

// Set stores a fixture credential in this recorder only.
func (r *Recorder) Set(service, user, password string) error {
	if err := r.record(Intent{Operation: "secret.save"}); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.values == nil {
		r.values = map[string]string{}
	}
	r.values[service+"\x00"+user] = password
	return nil
}

// Get loads an isolated fixture credential without logging its value.
func (r *Recorder) Get(service, user string) (string, error) {
	if err := r.record(Intent{Operation: "secret.load"}); err != nil {
		return "", err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.values[service+"\x00"+user]
	if !ok {
		return "", errors.New("fixture credential does not exist")
	}
	return value, nil
}

// Delete removes only a fixture credential.
func (r *Recorder) Delete(service, user string) error {
	if err := r.record(Intent{Operation: "secret.delete"}); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.values, service+"\x00"+user)
	return nil
}

// GetContext reads only this recorder and honors invocation cancellation.
func (r *Recorder) GetContext(ctx context.Context, service, user string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return r.Get(service, user)
}
