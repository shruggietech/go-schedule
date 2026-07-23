//go:build !windows

package service

import "github.com/kardianos/service"

// platformStatus has no platform-specific implementation outside Windows. It
// reports "not handled" so Control falls through to the library's own Status(),
// which is exactly what it did before the Windows least-privilege path existed
// — the Linux and macOS behavior is unchanged rather than reimplemented.
func platformStatus(string) (service.Status, bool, error) {
	return service.StatusUnknown, false, nil
}
