//go:build !windows

package gui

func diagnoseAccess() accessDiagnosis {
	return accessDiagnosis{Detail: "Windows group and token diagnostics are unavailable on this platform."}
}
