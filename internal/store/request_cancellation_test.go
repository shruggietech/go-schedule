package store_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/ipc"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestOwnedIPCRequestCancellationReleasesPoolWait(t *testing.T) {
	operations := []struct {
		name string
		run  func(context.Context, *client.Client) error
	}{
		{"tasks", func(ctx context.Context, c *client.Client) error { _, err := c.ListTasks(ctx, "", ""); return err }},
		{"runtime_info", func(ctx context.Context, c *client.Client) error { _, err := c.RuntimeInfo(ctx); return err }},
		{"credentials", func(ctx context.Context, c *client.Client) error { _, err := c.ListCredentials(ctx); return err }},
		{"group_delete", func(ctx context.Context, c *client.Client) error { return c.DeleteGroup(ctx, "owned-target") }},
	}
	for _, operation := range operations {
		for _, database := range []string{":memory:", filepath.Join(t.TempDir(), "owned.db")} {
			t.Run(fmt.Sprintf("%s/%s", operation.name, filepath.Base(database)), func(t *testing.T) {
				st, err := store.Open(database)
				if err != nil {
					t.Fatal(err)
				}
				defer st.Close()
				w, err := clifixture.NewWorkspace("")
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := w.CleanupRuntime(); err != nil {
						t.Error(err)
					}
				}()
				listener, _, err := ipc.Listen(w.Config)
				if err != nil {
					t.Fatal(err)
				}
				api := server.New(st, nil, nil, nil, "", nil)
				finished := make(chan struct{}, 1)
				svc := &http.Server{ReadHeaderTimeout: time.Second, Handler: http.HandlerFunc(func(out http.ResponseWriter, r *http.Request) {
					api.Handler().ServeHTTP(out, r)
					finished <- struct{}{}
				})}
				serveDone := make(chan error, 1)
				go func() { serveDone <- svc.Serve(listener) }()
				defer func() {
					_ = svc.Close()
					if err := <-serveDone; !errors.Is(err, http.ErrServerClosed) {
						t.Error(err)
					}
				}()
				held, err := st.HoldConnectionForTest(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				defer held.Close()
				baseline := st.PoolStatsForTest().WaitCount
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				result := make(chan error, 1)
				go func() { result <- operation.run(ctx, client.New(w.IPCPath)) }()
				deadline := time.NewTimer(5 * time.Second)
				defer deadline.Stop()
				tick := time.NewTicker(time.Millisecond)
				defer tick.Stop()
				for st.PoolStatsForTest().WaitCount == baseline {
					select {
					case <-tick.C:
					case <-deadline.C:
						t.Fatal("request never entered owned SQL pool wait")
					}
				}
				t.Logf("queued real IPC request: %+v", st.PoolStatsForTest())
				cancel()
				select {
				case err := <-result:
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("client cancellation: %v", err)
					}
				case <-time.After(5 * time.Second):
					t.Fatal("client failed to cancel")
				}
				select {
				case <-finished:
				case <-time.After(time.Second):
					t.Error("server handler remained blocked on owned SQL pool after client cancellation")
				}
				// Release only after observing server completion, including on a failed regression.
				_ = held.Close()
			})
		}
	}
}
