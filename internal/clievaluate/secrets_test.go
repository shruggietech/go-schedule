package clievaluate

import (
	"bytes"
	"strings"
	"testing"
)

func TestIdentityKeysRemainVisibleAndTriggerKeysAreRedacted(t *testing.T) {
	redactor := NewRedactor()
	identity := []byte(`{"key":"public-profile-identifier","phrase":"weekdays at 09:00"}`)
	if actual := redactor.Sanitize(identity); !bytes.Equal(actual, identity) {
		t.Fatal("public identity or schedule phrase was treated as a credential")
	}
	description := []byte("A bearer credential is required to connect.")
	if !bytes.Equal(redactor.Sanitize(description), description) {
		t.Fatal("ordinary help text was treated as an authorization header")
	}
	registerTriggerKeys(redactor, map[string]any{"members": []any{map[string]any{"key": "actual-trigger-secret"}}})
	actual := redactor.Sanitize([]byte("gosched trigger fire actual-trigger-secret"))
	if strings.Contains(string(actual), "actual-trigger-secret") {
		t.Fatal("registered trigger key survived redaction")
	}
}

func TestHumanBindingsRejectAmbiguousIdentityAndExtractOneTimeLabels(t *testing.T) {
	before := map[string]any{"groups": map[string]any{"existing": nil}}
	after := map[string]any{"groups": map[string]any{"existing": nil, "created": nil}}
	id, err := newIdentity(before, after, "groups")
	if err != nil || id != "created" {
		t.Fatal(id, err)
	}
	after["groups"].(map[string]any)["second"] = nil
	if _, err := newIdentity(before, after, "groups"); err == nil {
		t.Fatal("ambiguous creation was accepted")
	}
	if actual := humanLine("Pairing ID: fixture\nPhrase (shown once): owned phrase\n", "Phrase (shown once):"); actual != "owned phrase" {
		t.Fatalf("phrase label=%q", actual)
	}
}
