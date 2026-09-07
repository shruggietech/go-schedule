package connection

import (
	"context"
	"sync"
	"time"
)

type backendResult struct {
	health Health
	err    error
}

type backendFake struct {
	mu           sync.Mutex
	results      []backendResult
	streams      []chan DomainEvent
	streamErrors []chan error
}

func (f *backendFake) Health(context.Context) (Health, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.results) == 0 {
		return Health{}, &Failure{State: StateUnavailable, Message: "Unavailable.", Action: "Retry."}
	}
	result := f.results[0]
	if len(f.results) > 1 {
		f.results = f.results[1:]
	}
	return result.health, result.err
}

func (f *backendFake) StreamEvents(ctx context.Context, publish func(DomainEvent)) error {
	eventsCh := make(chan DomainEvent, 4)
	errorsCh := make(chan error, 1)
	f.mu.Lock()
	f.streams = append(f.streams, eventsCh)
	f.streamErrors = append(f.streamErrors, errorsCh)
	f.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case event := <-eventsCh:
		publish(event)
		<-ctx.Done()
		return ctx.Err()
	case err := <-errorsCh:
		return err
	}
}

type observerFake struct{ events chan Event }

func (o observerFake) Publish(event Event) { o.events <- event }

type schedulerFake struct {
	delays   chan time.Duration
	releases chan func()
}

func (s schedulerFake) After(ctx context.Context, delay time.Duration) <-chan struct{} {
	ready := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(ready) }) }
	s.delays <- delay
	s.releases <- release
	go func() {
		<-ctx.Done()
		release()
	}()
	return ready
}
