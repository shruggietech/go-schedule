package clievaluate

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestSMTPSecretsStayOutsideCommandSubstitutions(t *testing.T) {
	w, err := clifixture.NewWorkspace("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(w.Root); err != nil {
			t.Error(err)
		}
	})
	s := &session{w: w, vars: map[string]string{}, redactor: NewRedactor()}
	if err := s.startSMTP(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := s.smtp.Close(); err != nil {
			t.Error(err)
		}
	})
	data, err := os.ReadFile(s.vars["smtp_auth_file"])
	if err != nil {
		t.Fatal(err)
	}
	var auth domain.SMTPAuth
	if err := json.Unmarshal(data, &auth); err != nil {
		t.Fatal(err)
	}
	if auth.Username == "" || auth.Password == "" || s.smtpRotation.Password == "" || s.smtpRotation.Password == auth.Password {
		t.Fatal("missing independent credentials")
	}
	for key, value := range s.vars {
		if value == auth.Username || value == auth.Password || value == s.smtpRotation.Password {
			t.Fatalf("credential entered generic substitution %s", key)
		}
	}
	for _, secret := range []string{auth.Username, auth.Password, s.smtpRotation.Password} {
		if got := string(s.redactor.Sanitize([]byte(secret))); got != "[REDACTED]" {
			t.Fatal("credential was not registered for redaction")
		}
	}
}
