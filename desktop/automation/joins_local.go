package automation

import (
	"context"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func (b *LocalBackend) CreateJoin(c context.Context, r server.JoinCreateRequest) (domain.ParallelJoin, error) {
	return b.daemon.CreateJoin(c, r)
}
func (b *LocalBackend) UpdateJoin(c context.Context, id string, r server.JoinUpdateRequest) (domain.ParallelJoin, error) {
	return b.daemon.UpdateJoin(c, id, r)
}
func (b *LocalBackend) DeleteJoin(c context.Context, id string) error {
	return b.daemon.DeleteJoin(c, id)
}
func (b *LocalBackend) ListJoins(c context.Context) ([]domain.ParallelJoin, error) {
	return b.daemon.ListJoins(c)
}
func (b *LocalBackend) GetJoin(c context.Context, id string) (domain.ParallelJoin, error) {
	return b.daemon.GetJoin(c, id)
}
func (b *LocalBackend) ListJoinInstances(c context.Context, f domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, error) {
	return b.daemon.ListJoinInstances(c, f)
}
func (b *LocalBackend) GetJoinInstance(c context.Context, id string) (domain.ParallelJoinInstance, error) {
	return b.daemon.GetJoinInstance(c, id)
}
