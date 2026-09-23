package service

// State is the exact lifecycle state of the installed system service.
// Unknown means the manager could not classify it; NotInstalled is distinct.
type State string

const (
	StateNotInstalled State = "not_installed"
	StateStopped      State = "stopped"
	StateStarting     State = "starting"
	StateRunning      State = "running"
	StateStopping     State = "stopping"
	StateUnknown      State = "unknown"
)

// QueryState reads the local installed service lifecycle without requesting
// permission to start or stop it.
func QueryState() (State, error) {
	return platformQueryState(svcName)
}
