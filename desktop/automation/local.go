package automation

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Backend interface {
	ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error)
	ListChains(context.Context) ([]domain.CompletionChain, error)
	GetChain(context.Context, string) (domain.CompletionChain, error)
	CreateChain(context.Context, server.ChainCreateRequest) (domain.CompletionChain, error)
	UpdateChain(context.Context, string, server.ChainUpdateRequest) (domain.CompletionChain, error)
	DeleteChain(context.Context, string) error
	ListTriggers(context.Context) ([]server.TriggerResponse, error)
	GetTrigger(context.Context, string) (server.TriggerResponse, error)
	CreateTrigger(context.Context, server.TriggerCreateRequest) (server.TriggerSecretResponse, error)
	UpdateTrigger(context.Context, string, server.TriggerUpdateRequest) (server.TriggerResponse, error)
	SetTriggerEnabled(context.Context, string, bool) (server.TriggerResponse, error)
	RevealTrigger(context.Context, string) (server.TriggerSecretResponse, error)
	RotateTrigger(context.Context, string) (server.TriggerSecretResponse, error)
	FireTrigger(context.Context, string) error
	DeleteTrigger(context.Context, string) error
	ListTriggerSets(context.Context) ([]server.TriggerSetResponse, error)
	GetTriggerSet(context.Context, string) (server.TriggerSetResponse, error)
	CreateTriggerSet(context.Context, server.TriggerSetCreateRequest) (server.TriggerSetSecretResponse, error)
	RetargetTriggerSet(context.Context, string, string) (server.TriggerSetResponse, error)
	SetTriggerSetEnabled(context.Context, string, bool) (server.TriggerSetResponse, error)
	RevealTriggerSet(context.Context, string) (server.TriggerSetSecretResponse, error)
	RotateTriggerSet(context.Context, string) (server.TriggerSetSecretResponse, error)
	DeleteTriggerSet(context.Context, string) error
	ListFilesystemWatchers(context.Context) ([]server.FilesystemWatcherResponse, error)
	GetFilesystemWatcher(context.Context, string) (server.FilesystemWatcherResponse, error)
	CreateFilesystemWatcher(context.Context, server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error)
	UpdateFilesystemWatcher(context.Context, string, server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error)
	SetFilesystemWatcherEnabled(context.Context, string, bool) (server.FilesystemWatcherResponse, error)
	DeleteFilesystemWatcher(context.Context, string) error
}

type LocalBackend struct{ daemon Backend }

func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) ListTaskDetails(c context.Context, g, s string) ([]server.TaskResponse, error) {
	return b.daemon.ListTaskDetails(c, g, s)
}
func (b *LocalBackend) ListChains(c context.Context) ([]domain.CompletionChain, error) {
	return b.daemon.ListChains(c)
}
func (b *LocalBackend) GetChain(c context.Context, id string) (domain.CompletionChain, error) {
	return b.daemon.GetChain(c, id)
}
func (b *LocalBackend) CreateChain(c context.Context, r server.ChainCreateRequest) (domain.CompletionChain, error) {
	return b.daemon.CreateChain(c, r)
}
func (b *LocalBackend) UpdateChain(c context.Context, id string, r server.ChainUpdateRequest) (domain.CompletionChain, error) {
	return b.daemon.UpdateChain(c, id, r)
}
func (b *LocalBackend) DeleteChain(c context.Context, id string) error {
	return b.daemon.DeleteChain(c, id)
}
func (b *LocalBackend) ListTriggers(c context.Context) ([]server.TriggerResponse, error) {
	return b.daemon.ListTriggers(c)
}
func (b *LocalBackend) GetTrigger(c context.Context, id string) (server.TriggerResponse, error) {
	return b.daemon.GetTrigger(c, id)
}
func (b *LocalBackend) CreateTrigger(c context.Context, r server.TriggerCreateRequest) (server.TriggerSecretResponse, error) {
	return b.daemon.CreateTrigger(c, r)
}
func (b *LocalBackend) UpdateTrigger(c context.Context, id string, r server.TriggerUpdateRequest) (server.TriggerResponse, error) {
	return b.daemon.UpdateTrigger(c, id, r)
}
func (b *LocalBackend) SetTriggerEnabled(c context.Context, id string, e bool) (server.TriggerResponse, error) {
	return b.daemon.SetTriggerEnabled(c, id, e)
}
func (b *LocalBackend) RevealTrigger(c context.Context, id string) (server.TriggerSecretResponse, error) {
	return b.daemon.RevealTrigger(c, id)
}
func (b *LocalBackend) RotateTrigger(c context.Context, id string) (server.TriggerSecretResponse, error) {
	return b.daemon.RotateTrigger(c, id)
}
func (b *LocalBackend) FireTrigger(c context.Context, key string) error {
	return b.daemon.FireTrigger(c, key)
}
func (b *LocalBackend) DeleteTrigger(c context.Context, id string) error {
	return b.daemon.DeleteTrigger(c, id)
}
func (b *LocalBackend) ListTriggerSets(c context.Context) ([]server.TriggerSetResponse, error) {
	return b.daemon.ListTriggerSets(c)
}
func (b *LocalBackend) GetTriggerSet(c context.Context, id string) (server.TriggerSetResponse, error) {
	return b.daemon.GetTriggerSet(c, id)
}
func (b *LocalBackend) CreateTriggerSet(c context.Context, r server.TriggerSetCreateRequest) (server.TriggerSetSecretResponse, error) {
	return b.daemon.CreateTriggerSet(c, r)
}
func (b *LocalBackend) RetargetTriggerSet(c context.Context, id, task string) (server.TriggerSetResponse, error) {
	return b.daemon.RetargetTriggerSet(c, id, task)
}
func (b *LocalBackend) SetTriggerSetEnabled(c context.Context, id string, e bool) (server.TriggerSetResponse, error) {
	return b.daemon.SetTriggerSetEnabled(c, id, e)
}
func (b *LocalBackend) RevealTriggerSet(c context.Context, id string) (server.TriggerSetSecretResponse, error) {
	return b.daemon.RevealTriggerSet(c, id)
}
func (b *LocalBackend) RotateTriggerSet(c context.Context, id string) (server.TriggerSetSecretResponse, error) {
	return b.daemon.RotateTriggerSet(c, id)
}
func (b *LocalBackend) DeleteTriggerSet(c context.Context, id string) error {
	return b.daemon.DeleteTriggerSet(c, id)
}
func (b *LocalBackend) ListFilesystemWatchers(c context.Context) ([]server.FilesystemWatcherResponse, error) {
	return b.daemon.ListFilesystemWatchers(c)
}
func (b *LocalBackend) GetFilesystemWatcher(c context.Context, id string) (server.FilesystemWatcherResponse, error) {
	return b.daemon.GetFilesystemWatcher(c, id)
}
func (b *LocalBackend) CreateFilesystemWatcher(c context.Context, r server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error) {
	return b.daemon.CreateFilesystemWatcher(c, r)
}
func (b *LocalBackend) UpdateFilesystemWatcher(c context.Context, id string, r server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error) {
	return b.daemon.UpdateFilesystemWatcher(c, id, r)
}
func (b *LocalBackend) SetFilesystemWatcherEnabled(c context.Context, id string, e bool) (server.FilesystemWatcherResponse, error) {
	return b.daemon.SetFilesystemWatcherEnabled(c, id, e)
}
func (b *LocalBackend) DeleteFilesystemWatcher(c context.Context, id string) error {
	return b.daemon.DeleteFilesystemWatcher(c, id)
}
