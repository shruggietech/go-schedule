package notification

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/shruggietech/go-schedule/internal/clock"
	"github.com/shruggietech/go-schedule/internal/domain"
)

const (
	maxAttempts  = 3
	pollInterval = 250 * time.Millisecond
)

// DeliveryStore is the durable state used by Dispatcher.
type DeliveryStore interface {
	RecoverNotificationDeliveries(time.Time) (int64, error)
	ClaimNotificationDeliveries(int, time.Time) ([]domain.NotificationDelivery, error)
	RetryNotificationDelivery(string, time.Time, int, string) error
	CompleteNotificationDelivery(string, bool, int, string, time.Time) error
}

// AttemptSender performs one outbound attempt.
type AttemptSender interface {
	Send(context.Context, domain.NotificationDelivery) (int, string, error)
}

// Dispatcher owns the bounded notification worker lifecycle. Wake is
// coalesced, Run owns all attempt goroutines, and Run returns only after they end.
type Dispatcher struct {
	store   DeliveryStore
	sender  AttemptSender
	log     *slog.Logger
	clock   clock.Clock
	workers int
	wake    chan struct{}
}

// NewDispatcher creates a dispatcher with a dedicated outbound worker limit.
func NewDispatcher(store DeliveryStore, sender AttemptSender, log *slog.Logger, workers int) *Dispatcher {
	return NewDispatcherWithClock(store, sender, log, workers, clock.NewReal())
}

// NewDispatcherWithClock creates a dispatcher with an injected timing source.
func NewDispatcherWithClock(store DeliveryStore, sender AttemptSender, log *slog.Logger, workers int, clk clock.Clock) *Dispatcher {
	if workers <= 0 {
		workers = 4
	}
	if log == nil {
		log = slog.Default()
	}
	return &Dispatcher{store: store, sender: sender, log: log, clock: clk, workers: workers, wake: make(chan struct{}, 1)}
}

// Wake requests a prompt scan for newly committed work.
func (d *Dispatcher) Wake() {
	select {
	case d.wake <- struct{}{}:
	default:
	}
}

// Run recovers interrupted work and dispatches until context cancellation.
func (d *Dispatcher) Run(ctx context.Context) error {
	now := d.clock.Now().UTC()
	if recovered, err := d.store.RecoverNotificationDeliveries(now); err != nil {
		return err
	} else if recovered > 0 {
		d.log.Warn("notification: recovered interrupted deliveries", "count", recovered)
	}
	for {
		if err := d.process(ctx); err != nil {
			d.log.Error("notification: process deliveries", "err", err)
		}
		timer := d.clock.NewTimer(pollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-d.wake:
			timer.Stop()
		case <-timer.C:
		}
	}
}

func (d *Dispatcher) process(ctx context.Context) error {
	deliveries, err := d.store.ClaimNotificationDeliveries(d.workers, d.clock.Now().UTC())
	if err != nil || len(deliveries) == 0 {
		return err
	}
	var wg sync.WaitGroup
	for _, delivery := range deliveries {
		delivery := delivery
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, diagnostic, sendErr := d.sender.Send(ctx, delivery)
			d.finishAttempt(delivery, status, diagnostic, sendErr)
		}()
	}
	wg.Wait()
	return nil
}

func (d *Dispatcher) finishAttempt(delivery domain.NotificationDelivery, status int, diagnostic string, sendErr error) {
	now := d.clock.Now().UTC()
	if sendErr == nil {
		if err := d.store.CompleteNotificationDelivery(delivery.ID, true, status, "", now); err != nil {
			d.log.Error("notification: record success", "delivery", delivery.ID, "channel", delivery.ChannelID, "err", err)
		}
		return
	}
	if delivery.Attempts < maxAttempts {
		delay := time.Duration(1<<(delivery.Attempts-1)) * time.Second
		if err := d.store.RetryNotificationDelivery(delivery.ID, now.Add(delay), status, diagnostic); err != nil {
			d.log.Error("notification: schedule retry", "delivery", delivery.ID, "channel", delivery.ChannelID, "attempt", delivery.Attempts, "err", err)
		}
		return
	}
	if err := d.store.CompleteNotificationDelivery(delivery.ID, false, status, diagnostic, now); err != nil {
		d.log.Error("notification: record failure", "delivery", delivery.ID, "channel", delivery.ChannelID, "attempts", delivery.Attempts, "err", err)
		return
	}
	d.log.Warn("notification: delivery failed", "delivery", delivery.ID, "channel", delivery.ChannelID, "task", delivery.TaskID, "run", delivery.RunID, "attempts", delivery.Attempts, "status", status, "reason", diagnostic)
}
