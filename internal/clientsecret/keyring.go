package clientsecret

import (
	"errors"
	"strings"

	keyring "github.com/zalando/go-keyring"
)

const servicePrefix = "go-schedule.remote."

type Backend interface {
	Set(service, user, password string) error
	Get(service, user string) (string, error)
	Delete(service, user string) error
}

type nativeBackend struct{}

func (nativeBackend) Set(service, user, password string) error {
	return keyring.Set(service, user, password)
}
func (nativeBackend) Get(service, user string) (string, error) { return keyring.Get(service, user) }
func (nativeBackend) Delete(service, user string) error        { return keyring.Delete(service, user) }

type Store struct{ backend Backend }

func New() *Store                           { return &Store{backend: nativeBackend{}} }
func NewWithBackend(backend Backend) *Store { return &Store{backend: backend} }

func (s *Store) Probe() error {
	const service, user, value = "go-schedule.remote.probe", "availability", "probe"
	if err := s.backend.Set(service, user, value); err != nil {
		return err
	}
	read, err := s.backend.Get(service, user)
	deleteErr := s.backend.Delete(service, user)
	if err != nil {
		return err
	}
	if deleteErr != nil {
		return deleteErr
	}
	if read != value {
		return errors.New("native credential store verification failed")
	}
	return nil
}

func (s *Store) Save(daemonID, credentialID, token string) error {
	if strings.TrimSpace(daemonID) == "" || strings.TrimSpace(credentialID) == "" || token == "" {
		return errors.New("credential identity is incomplete")
	}
	return s.backend.Set(servicePrefix+daemonID, credentialID, token)
}

func (s *Store) Load(daemonID, credentialID string) (string, error) {
	return s.backend.Get(servicePrefix+daemonID, credentialID)
}
func (s *Store) Delete(daemonID, credentialID string) error {
	return s.backend.Delete(servicePrefix+daemonID, credentialID)
}
