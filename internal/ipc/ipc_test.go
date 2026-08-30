package ipc

import "testing"

func TestAccessModeContract(t *testing.T) {
	if AccessModeRestricted != "restricted" {
		t.Fatalf("restricted access mode = %q", AccessModeRestricted)
	}
	if AccessModeCompatibility != "compatibility" {
		t.Fatalf("compatibility access mode = %q", AccessModeCompatibility)
	}

	restricted := AccessInfo{Mode: AccessModeRestricted, AdminGroup: "goschedadmin"}
	if restricted.Mode != AccessModeRestricted || restricted.AdminGroup == "" {
		t.Fatalf("restricted access evidence = %+v", restricted)
	}
	compatibility := AccessInfo{Mode: AccessModeCompatibility}
	if compatibility.Mode != AccessModeCompatibility || compatibility.AdminGroup != "" {
		t.Fatalf("compatibility access evidence = %+v", compatibility)
	}
}
