package operations

import (
	"context"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type Backend interface {
	GetCalendar(context.Context, time.Time, time.Time) (server.CalendarResponse, error)
	ListRuns(context.Context, string, int) ([]domain.Run, error)
	ListLogs(context.Context, string, int) (server.LogsResponse, error)
	ListAlerts(context.Context, bool) ([]domain.Alert, error)
	AckAlert(context.Context, string) error
}

type LocalBackend struct{ daemon Backend }

func NewLocalBackend(daemon Backend) *LocalBackend { return &LocalBackend{daemon: daemon} }

func (b *LocalBackend) GetCalendar(ctx context.Context, from, to time.Time) (server.CalendarResponse, error) {
	return b.daemon.GetCalendar(ctx, from, to)
}
func (b *LocalBackend) ListRuns(ctx context.Context, task string, limit int) ([]domain.Run, error) {
	return b.daemon.ListRuns(ctx, task, limit)
}
func (b *LocalBackend) ListLogs(ctx context.Context, severity string, limit int) (server.LogsResponse, error) {
	return b.daemon.ListLogs(ctx, severity, limit)
}
func (b *LocalBackend) ListAlerts(ctx context.Context, unacked bool) ([]domain.Alert, error) {
	return b.daemon.ListAlerts(ctx, unacked)
}
func (b *LocalBackend) AckAlert(ctx context.Context, id string) error {
	return b.daemon.AckAlert(ctx, id)
}
