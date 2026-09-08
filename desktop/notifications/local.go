package notifications

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Backend interface {
	ListNotificationChannels(context.Context) ([]domain.NotificationChannel, error)
	CreateNotificationChannel(context.Context, server.NotificationChannelCreateRequest) (domain.NotificationChannel, error)
	UpdateNotificationChannel(context.Context, string, server.NotificationChannelUpdateRequest) (domain.NotificationChannel, error)
	SetNotificationChannelEnabled(context.Context, string, bool) (domain.NotificationChannel, error)
	RotateNotificationChannelAuthorization(context.Context, string, string) (domain.NotificationChannel, error)
	TestNotificationChannel(context.Context, string) (domain.NotificationDelivery, error)
	DeleteNotificationChannel(context.Context, string) error
	ListNotificationDeliveries(context.Context, domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error)
	ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error)
	ListGroups(context.Context) ([]domain.Group, error)
	ListTaskNotificationAssignments(context.Context, string) ([]domain.NotificationAssignment, error)
	ListGroupNotificationAssignments(context.Context, string) ([]domain.NotificationAssignment, error)
	EffectiveTaskNotificationPolicy(context.Context, string) (domain.EffectiveNotificationPolicy, error)
	ReplaceTaskNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error)
	ReplaceGroupNotificationAssignments(context.Context, string, server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error)
}

type LocalBackend struct{ daemon Backend }

func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) ListNotificationChannels(c context.Context) ([]domain.NotificationChannel, error) {
	return b.daemon.ListNotificationChannels(c)
}
func (b *LocalBackend) CreateNotificationChannel(c context.Context, r server.NotificationChannelCreateRequest) (domain.NotificationChannel, error) {
	return b.daemon.CreateNotificationChannel(c, r)
}
func (b *LocalBackend) UpdateNotificationChannel(c context.Context, id string, r server.NotificationChannelUpdateRequest) (domain.NotificationChannel, error) {
	return b.daemon.UpdateNotificationChannel(c, id, r)
}
func (b *LocalBackend) SetNotificationChannelEnabled(c context.Context, id string, enabled bool) (domain.NotificationChannel, error) {
	return b.daemon.SetNotificationChannelEnabled(c, id, enabled)
}
func (b *LocalBackend) RotateNotificationChannelAuthorization(c context.Context, id, value string) (domain.NotificationChannel, error) {
	return b.daemon.RotateNotificationChannelAuthorization(c, id, value)
}
func (b *LocalBackend) TestNotificationChannel(c context.Context, id string) (domain.NotificationDelivery, error) {
	return b.daemon.TestNotificationChannel(c, id)
}
func (b *LocalBackend) DeleteNotificationChannel(c context.Context, id string) error {
	return b.daemon.DeleteNotificationChannel(c, id)
}
func (b *LocalBackend) ListNotificationDeliveries(c context.Context, f domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error) {
	return b.daemon.ListNotificationDeliveries(c, f)
}
func (b *LocalBackend) ListTaskDetails(c context.Context, g, s string) ([]server.TaskResponse, error) {
	return b.daemon.ListTaskDetails(c, g, s)
}
func (b *LocalBackend) ListGroups(c context.Context) ([]domain.Group, error) {
	return b.daemon.ListGroups(c)
}
func (b *LocalBackend) ListTaskNotificationAssignments(c context.Context, id string) ([]domain.NotificationAssignment, error) {
	return b.daemon.ListTaskNotificationAssignments(c, id)
}
func (b *LocalBackend) ListGroupNotificationAssignments(c context.Context, id string) ([]domain.NotificationAssignment, error) {
	return b.daemon.ListGroupNotificationAssignments(c, id)
}
func (b *LocalBackend) EffectiveTaskNotificationPolicy(c context.Context, id string) (domain.EffectiveNotificationPolicy, error) {
	return b.daemon.EffectiveTaskNotificationPolicy(c, id)
}
func (b *LocalBackend) ReplaceTaskNotificationAssignments(c context.Context, id string, r server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	return b.daemon.ReplaceTaskNotificationAssignments(c, id, r)
}
func (b *LocalBackend) ReplaceGroupNotificationAssignments(c context.Context, id string, r server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	return b.daemon.ReplaceGroupNotificationAssignments(c, id, r)
}
