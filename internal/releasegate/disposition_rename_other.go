//go:build !darwin && !linux && !windows

package releasegate

import "fmt"

func renameDispositionNoReplace(_, _ string) error {
	return fmt.Errorf("atomic no-replace disposition commit is unsupported on this platform")
}
