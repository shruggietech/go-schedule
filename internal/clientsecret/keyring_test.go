package clientsecret

import "testing"

type fakeBackend struct {
	values map[string]string
	fail   bool
}

func (f *fakeBackend) Set(s, u, p string) error {
	if f.fail {
		return errFake
	}
	f.values[s+u] = p
	return nil
}
func (f *fakeBackend) Get(s, u string) (string, error) {
	v, ok := f.values[s+u]
	if !ok {
		return "", errFake
	}
	return v, nil
}
func (f *fakeBackend) Delete(s, u string) error { delete(f.values, s+u); return nil }

type fakeError string

func (e fakeError) Error() string { return string(e) }

const errFake = fakeError("keyring unavailable")

func TestStoreUsesDaemonScopedNativeCredentialEntry(t *testing.T) {
	backend := &fakeBackend{values: map[string]string{}}
	s := NewWithBackend(backend)
	if err := s.Probe(); err != nil {
		t.Fatal(err)
	}
	if err := s.Save("daemon", "credential", "secret"); err != nil {
		t.Fatal(err)
	}
	value, err := s.Load("daemon", "credential")
	if err != nil || value != "secret" {
		t.Fatalf("value=%q err=%v", value, err)
	}
}
func TestProbeFailsClosed(t *testing.T) {
	if NewWithBackend(&fakeBackend{values: map[string]string{}, fail: true}).Probe() == nil {
		t.Fatal("expected unavailable keyring")
	}
}
