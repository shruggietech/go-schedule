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
// requested mask at once — so a read-only status query failed with "Access is
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
	scm, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return service.StatusUnknown, true, err
	}
	defer windows.CloseServiceHandle(scm) //nolint:errcheck // read-only handle

	svcName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return service.StatusUnknown, true, err
	}

	h, err := windows.OpenService(scm, svcName, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		if errors.Is(err, windows.ERROR_SERVICE_DOES_NOT_EXIST) {
			// Same result the library reports for a missing service, so the CLI
			// prints the existing "unknown (is the service installed?)" wording.
			return service.StatusUnknown, true, nil
		}
		return service.StatusUnknown, true, err
	}
	defer windows.CloseServiceHandle(h) //nolint:errcheck // read-only handle

	var st windows.SERVICE_STATUS
	if err := windows.QueryServiceStatus(h, &st); err != nil {
		return service.StatusUnknown, true, err
	}

	switch st.CurrentState {
	case windows.SERVICE_RUNNING, windows.SERVICE_START_PENDING:
		return service.StatusRunning, true, nil
	case windows.SERVICE_STOPPED, windows.SERVICE_STOP_PENDING,
		windows.SERVICE_PAUSED, windows.SERVICE_PAUSE_PENDING,
		windows.SERVICE_CONTINUE_PENDING:
		return service.StatusStopped, true, nil
	default:
		return service.StatusUnknown, true, nil
	}
}
