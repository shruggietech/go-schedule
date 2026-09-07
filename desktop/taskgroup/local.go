package taskgroup

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Backend interface {
	ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error)
	GetTask(context.Context, string) (server.TaskResponse, error)
	CreateTask(context.Context, server.TaskCreateRequest) (server.TaskResponse, error)
	UpdateTask(context.Context, string, server.TaskUpdateRequest) (server.TaskResponse, error)
	DeleteTask(context.Context, string) error
	SetTaskEnabled(context.Context, string, bool) error
	RunNow(context.Context, string) error
	Preview(context.Context, server.PreviewRequest) (server.PreviewResponse, error)
	ListGroups(context.Context) ([]domain.Group, error)
	CreateGroup(context.Context, server.GroupCreateRequest) (domain.Group, error)
	UpdateGroup(context.Context, string, server.GroupUpdateRequest) (domain.Group, error)
	SetGroupEnabled(context.Context, string, bool) error
	DeleteGroup(context.Context, string) error
}

type LocalBackend struct{ daemon Backend }

func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) ListTaskDetails(c context.Context, g, s string) ([]server.TaskResponse, error) {
	return b.daemon.ListTaskDetails(c, g, s)
}
func (b *LocalBackend) GetTask(c context.Context, id string) (server.TaskResponse, error) {
	return b.daemon.GetTask(c, id)
}
func (b *LocalBackend) CreateTask(c context.Context, r server.TaskCreateRequest) (server.TaskResponse, error) {
	return b.daemon.CreateTask(c, r)
}
func (b *LocalBackend) UpdateTask(c context.Context, id string, r server.TaskUpdateRequest) (server.TaskResponse, error) {
	return b.daemon.UpdateTask(c, id, r)
}
func (b *LocalBackend) DeleteTask(c context.Context, id string) error {
	return b.daemon.DeleteTask(c, id)
}
func (b *LocalBackend) SetTaskEnabled(c context.Context, id string, e bool) error {
	return b.daemon.SetTaskEnabled(c, id, e)
}
func (b *LocalBackend) RunNow(c context.Context, id string) error { return b.daemon.RunNow(c, id) }
func (b *LocalBackend) Preview(c context.Context, r server.PreviewRequest) (server.PreviewResponse, error) {
	return b.daemon.Preview(c, r)
}
func (b *LocalBackend) ListGroups(c context.Context) ([]domain.Group, error) {
	return b.daemon.ListGroups(c)
}
func (b *LocalBackend) CreateGroup(c context.Context, r server.GroupCreateRequest) (domain.Group, error) {
	return b.daemon.CreateGroup(c, r)
}
func (b *LocalBackend) UpdateGroup(c context.Context, id string, r server.GroupUpdateRequest) (domain.Group, error) {
	return b.daemon.UpdateGroup(c, id, r)
}
func (b *LocalBackend) SetGroupEnabled(c context.Context, id string, e bool) error {
	return b.daemon.SetGroupEnabled(c, id, e)
}
func (b *LocalBackend) DeleteGroup(c context.Context, id string) error {
	return b.daemon.DeleteGroup(c, id)
}
