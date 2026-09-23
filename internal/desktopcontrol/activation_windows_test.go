//go:build windows

package desktopcontrol

import "testing"

func TestTrayInstanceIsSessionSingleton(t *testing.T) {
	const testName = "Local\\go-schedule-s101-instance-test"
	first, owner, err := claim(testName)
	if err != nil || !owner {
		t.Fatalf("first claim: owner=%v err=%v", owner, err)
	}
	defer first.Close() //nolint:errcheck // test cleanup
	second, owner, err := claim(testName)
	if err != nil || owner || second != nil {
		t.Fatalf("duplicate claim: instance=%v owner=%v err=%v", second, owner, err)
	}
}
