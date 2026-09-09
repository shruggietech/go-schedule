package mcpobserve

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestBoundedTextNormalizesAndPreservesUTF8Boundary(t *testing.T) {
	got, truncated := boundedText("ab\xff世界", 6)
	if !truncated {
		t.Fatal("boundedText did not report truncation")
	}
	if !utf8.ValidString(got) || len(got) > 6 {
		t.Fatalf("boundedText = %q (%d bytes)", got, len(got))
	}
	if !strings.Contains(got, "�") {
		t.Fatalf("boundedText did not normalize invalid UTF-8: %q", got)
	}
}

func TestCursorRoundTripAndValidation(t *testing.T) {
	encoded := encodeCursor(200)
	if got, err := decodeCursor(encoded); err != nil || got != 200 {
		t.Fatalf("decodeCursor() = %d, %v", got, err)
	}
	for _, value := range []string{"", "not-base64", "e30", encodeCursor(-1)} {
		if _, err := decodeCursor(value); err == nil {
			t.Fatalf("decodeCursor(%q) succeeded", value)
		}
	}
}
