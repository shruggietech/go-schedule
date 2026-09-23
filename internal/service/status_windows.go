//go:build windows

package service

import (
	"errors"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows"
)

// platformStatus answers a status query using the minimum access rights the
// query needs.
//
// The library's own Status() opens the service handle with
// SERVICE_QUERY_CONFIG|SERVICE_QUERY_STATUS|SERVICE_START|SERVICE_STOP. The
// installed service's ACL grants Interactive Users the query rights but
// deliberately withholds start and stop, and OpenService evaluates the whole
// requested mask at once, so a read-only status query failed with "Access is
// denied" for any non-elevated user, reporting that the ACL forbade something
// the ACL in fact permitted.
//
// Only the status path is reimplemented here. install/uninstall/start/stop keep
// using the library, where the broader mask is genuinely required.
//
// The second return reports whether this platform handled the query at all, so
// the caller can fall through to the library on platforms with no
// implementation. On Windows it is always true: a real failure is an error, not
// a fallback, because falling back would silently reintroduce the wide mask.
func platformStatus(name string) (service.Status, bool, error) {
	state, err := platformQueryState(name)
	if err != nil {
		return service.StatusUnknown, true, err
	}
	switch state {
	case StateRunning, StateStarting:
		return service.StatusRunning, true, nil
	case StateStopped, StateStopping:
		return service.StatusStopped, true, nil
	default:
		return service.StatusUnknown, true, nil
	}
}

func platformQueryState(name string) (State, error) {
	scm, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return StateUnknown, err
	}
	defer windows.CloseServiceHandle(scm) //nolint:errcheck // read-only handle

	svcName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return StateUnknown, err
	}

	h, err := windows.OpenService(scm, svcName, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		if errors.Is(err, windows.ERROR_SERVICE_DOES_NOT_EXIST) {
			return StateNotInstalled, nil
		}
		return StateUnknown, err
	}
	defer windows.CloseServiceHandle(h) //nolint:errcheck // read-only handle

	var st windows.SERVICE_STATUS
	if err := windows.QueryServiceStatus(h, &st); err != nil {
		return StateUnknown, err
	}

	return classifyWindowsState(st.CurrentState), nil
}

func classifyWindowsState(current uint32) State {
	switch current {
	case windows.SERVICE_RUNNING:
		return StateRunning
	case windows.SERVICE_START_PENDING, windows.SERVICE_CONTINUE_PENDING:
		return StateStarting
	case windows.SERVICE_STOPPED:
		return StateStopped
	case windows.SERVICE_STOP_PENDING:
		return StateStopping
	default:
		return StateUnknown
	}
}
