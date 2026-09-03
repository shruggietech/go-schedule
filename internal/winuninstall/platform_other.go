//go:build !windows

package winuninstall

// Wipe reports unsupported execution outside Windows; the MSI helper is Windows-only.
func Wipe() Result {
	return Result{Schema: resultSchema, State: StateInternalError, Error: "Windows cleanup is unavailable on this platform"}
}
