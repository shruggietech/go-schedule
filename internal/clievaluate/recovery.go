package clievaluate

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clock"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/notification"
	"github.com/shruggietech/go-schedule/v2/internal/store"
	"github.com/shruggietech/go-schedule/v2/internal/watcher"
)

type deliverySender struct {
	mu    sync.Mutex
	calls int
}

func (s *deliverySender) Send(context.Context, domain.NotificationDelivery) (int, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if s.calls == 1 {
		return 503, "receiver returned HTTP 503", errors.New("fixture receiver unavailable")
	}
	return 204, "", nil
}

func notificationRetry(ctx context.Context) execution {
	st, err := store.Open(":memory:")
	if err != nil {
		return policyResult(nil, false, err)
	}
	defer func() { _ = st.Close() }()
	channel := domain.NotificationChannel{Name: "QA retry", Kind: domain.NotificationChannelWebhook, Endpoint: "http://127.0.0.1:1/owned", EndpointSummary: "http://127.0.0.1:1", Enabled: true}
	if err := st.CreateNotificationChannel(&channel); err != nil {
		return policyResult(nil, false, err)
	}
	delivery, err := st.CreateTestNotificationDelivery(channel.ID)
	if err != nil {
		return policyResult(nil, false, err)
	}
	clk := clock.NewFake(time.Now().UTC())
	sender := &deliverySender{}
	var logs bytes.Buffer
	dispatcher := notification.NewDispatcherWithClock(st, sender, slog.New(slog.NewTextHandler(&logs, nil)), 1, clk)
	child, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- dispatcher.Run(child) }()
	var first, final domain.NotificationDelivery
	err = waitUntil(ctx, func() bool {
		items, e := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: channel.ID, Limit: 10})
		if e == nil && len(items) == 1 && items[0].State == domain.NotificationDeliveryPending && items[0].Attempts == 1 {
			first = items[0]
			return true
		}
		return false
	})
	if err == nil {
		clk.Advance(2 * time.Second)
		dispatcher.Wake()
		err = waitUntil(ctx, func() bool {
			items, e := st.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: channel.ID, Limit: 10})
			if e == nil && len(items) == 1 && items[0].State == domain.NotificationDeliverySucceeded {
				final = items[0]
				return true
			}
			return false
		})
	}
	cancel()
	stopErr := <-done
	if errors.Is(stopErr, context.Canceled) {
		stopErr = nil
	}
	result := policyResult(map[string]any{"delivery_id": delivery.ID, "retry": first, "completed": final}, first.Attempts == 1 && final.Attempts == 2 && final.ID == delivery.ID, errors.Join(err, stopErr))
	result.Logs = logs.Bytes()
	return result
}

type watcherDefinitions struct{ items []domain.FilesystemWatcher }

func (w watcherDefinitions) ListFilesystemWatchers() ([]domain.FilesystemWatcher, error) {
	return append([]domain.FilesystemWatcher{}, w.items...), nil
}

type watcherDispatch struct{ events chan string }

func (w watcherDispatch) FireFilesystemWatcher(id string) error {
	select {
	case w.events <- id:
		return nil
	default:
		return errors.New("fixture dispatch channel exceeded its bound")
	}
}

func (s *session) watcherRecovery(ctx context.Context) execution {
	root := filepath.Join(s.w.Runtime, "watchers", "recovery")
	if err := s.w.Contains(root); err != nil {
		return policyResult(nil, false, err)
	}
	clk := clock.NewFake(time.Now().UTC())
	definition := domain.FilesystemWatcher{ID: "qa-recovery", Name: "QA recovery", Kind: domain.WatcherDirectory, Path: root, Pattern: "*.txt", Debounce: 25 * time.Millisecond, Stability: 25 * time.Millisecond, TargetTaskID: s.vars["base_task"], Enabled: true}
	dispatch := watcherDispatch{make(chan string, 8)}
	var logs bytes.Buffer
	manager := watcher.New(watcherDefinitions{[]domain.FilesystemWatcher{definition}}, dispatch, clk, slog.New(slog.NewTextHandler(&logs, nil)))
	child, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- manager.Run(child) }()
	select {
	case <-manager.Ready():
	case <-ctx.Done():
		cancel()
		<-done
		return policyResult(nil, false, ctx.Err())
	}
	initial := manager.Health(definition.ID)
	err := os.Mkdir(root, 0700)
	if err == nil {
		clk.Advance(3 * time.Second)
		err = waitUntil(ctx, func() bool {
			health := manager.Health(definition.ID)
			if health.State == domain.WatcherActive {
				return true
			}
			clk.Advance(250 * time.Millisecond)
			return false
		})
	}
	if err == nil {
		err = s.writeOwned(filepath.Join(root, "ready.txt"), []byte("owned recovery event"))
	}
	if err == nil {
		err = waitUntil(ctx, func() bool {
			select {
			case <-dispatch.events:
				return true
			default:
				clk.Advance(50 * time.Millisecond)
				return false
			}
		})
	}
	final := manager.Health(definition.ID)
	cancel()
	stopErr := <-done
	result := policyResult(map[string]any{"initial_health": initial, "recovered_health": final}, initial.State != domain.WatcherActive && final.State == domain.WatcherActive, errors.Join(err, stopErr))
	result.Logs = logs.Bytes()
	return result
}
