package secretstore

import (
	"runtime"
	"testing"
)

func TestProtectRoundTrip(t *testing.T) {
	for _, value := range []string{"", "https://user:pass@example.test/hook?token=secret", "Bearer héllo 世界"} {
		protected, err := Protect(value)
		if err != nil {
			t.Fatal(err)
		}
		if runtime.GOOS == "windows" && value != "" && protected == value {
			t.Fatal("Windows secret was stored without DPAPI protection")
		}
		plain, err := Unprotect(protected)
		if err != nil || plain != value {
			t.Fatalf("round trip=%q err=%v, want %q", plain, err, value)
		}
	}
}
