//go:build !windows && !linux && !darwin

package clientsecret

func readUnixPeerPID(int) (int, error)          { return 0, errReadUnavailable }
func readProcessExecutable(int) (string, error) { return "", errReadUnavailable }
