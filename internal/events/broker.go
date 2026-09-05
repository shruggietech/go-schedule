// Package events provides a small in-process publish/subscribe broker used to
// stream run-state changes and new alerts to connected GUI clients over
// Server-Sent Events. Publishing is non-blocking: a slow subscriber drops
// events rather than stalling the engine.
package events

import (
	"sync"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Kind classifies an event.
type Kind string

const (
	KindRun        Kind = "run"
	KindAlert      Kind = "alert"
	KindLog        Kind = "log"
	KindTask       Kind = "task"
	KindGroup      Kind = "group"
	KindChain      Kind = "chain"
	KindTrigger    Kind = "trigger"
	KindTriggerSet Kind = "trigger_set"
	KindWatcher    Kind = "filesystem_watcher"
)

// Verb describes a change to an entity in a task/group event.
type Verb string

const (
	VerbCreated Verb = "created"
	VerbUpdated Verb = "updated"
	VerbDeleted Verb = "deleted"
)

// TaskEvent describes a task change. Task is nil for deletions (only ID is set).
type TaskEvent struct {
	Verb Verb         `json:"verb"`
	ID   string       `json:"id"`
	Task *domain.Task `json:"task,omitempty"`
}

// GroupEvent describes a group change. Group is nil for deletions (only ID set).
type GroupEvent struct {
	Verb  Verb          `json:"verb"`
	ID    string        `json:"id"`
	Group *domain.Group `json:"group,omitempty"`
}

// ChainEvent describes a completion-chain change. Chain is nil for deletions.
type ChainEvent struct {
	Verb  Verb                    `json:"verb"`
	ID    string                  `json:"id"`
	Chain *domain.CompletionChain `json:"chain,omitempty"`
}

// TriggerEvent describes an external-trigger change without carrying its key.
type TriggerEvent struct {
	Verb    Verb                    `json:"verb"`
	ID      string                  `json:"id"`
	Trigger *domain.ExternalTrigger `json:"trigger,omitempty"`
}

// TriggerSetEvent describes a broad Trigger Set change without carrying keys.
type TriggerSetEvent struct {
	Verb       Verb               `json:"verb"`
	ID         string             `json:"id"`
	TriggerSet *domain.TriggerSet `json:"trigger_set,omitempty"`
}

// WatcherEvent describes a filesystem watcher change or runtime health transition without exposing its path.
type WatcherEvent struct {
	Verb   Verb                  `json:"verb"`
	ID     string                `json:"id"`
	Name   string                `json:"name,omitempty"`
	Health *domain.WatcherHealth `json:"health,omitempty"`
}

// Event is a single notification delivered to subscribers.
type Event struct {
	Kind       Kind              `json:"kind"`
	Run        *domain.Run       `json:"run,omitempty"`
	Alert      *domain.Alert     `json:"alert,omitempty"`
	Log        *domain.LogRecord `json:"log,omitempty"`
	Task       *TaskEvent        `json:"task,omitempty"`
	Group      *GroupEvent       `json:"group,omitempty"`
	Chain      *ChainEvent       `json:"chain,omitempty"`
	Trigger    *TriggerEvent     `json:"trigger,omitempty"`
	TriggerSet *TriggerSetEvent  `json:"trigger_set,omitempty"`
	Watcher    *WatcherEvent     `json:"filesystem_watcher,omitempty"`
}

// Broker fans out events to all current subscribers.
type Broker struct {
	mu   sync.RWMutex
	subs map[chan Event]struct{}
}

// NewBroker creates an empty broker.
func NewBroker() *Broker {
	return &Broker{subs: make(map[chan Event]struct{})}
}

// Subscribe registers a new subscriber, returning its event channel and an
// unsubscribe function. The channel is buffered; the caller should drain it.
func (b *Broker) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 64)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, ch)
			b.mu.Unlock()
			close(ch)
		})
	}
	return ch, cancel
}

// Publish delivers e to all subscribers without blocking. Events for a
// subscriber whose buffer is full are dropped for that subscriber.
func (b *Broker) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}

// PublishRun is a convenience for run events.
func (b *Broker) PublishRun(r domain.Run) { b.Publish(Event{Kind: KindRun, Run: &r}) }

// PublishAlert is a convenience for alert events.
func (b *Broker) PublishAlert(a domain.Alert) { b.Publish(Event{Kind: KindAlert, Alert: &a}) }

// PublishLog is a convenience for log events (satisfies logbus.Publisher).
func (b *Broker) PublishLog(r domain.LogRecord) { b.Publish(Event{Kind: KindLog, Log: &r}) }

// PublishTask is a convenience for task-change events.
func (b *Broker) PublishTask(verb Verb, id string, t *domain.Task) {
	b.Publish(Event{Kind: KindTask, Task: &TaskEvent{Verb: verb, ID: id, Task: t}})
}

// PublishGroup is a convenience for group-change events.
func (b *Broker) PublishGroup(verb Verb, id string, g *domain.Group) {
	b.Publish(Event{Kind: KindGroup, Group: &GroupEvent{Verb: verb, ID: id, Group: g}})
}

// PublishChain is a convenience for completion-chain change events.
func (b *Broker) PublishChain(verb Verb, id string, chain *domain.CompletionChain) {
	b.Publish(Event{Kind: KindChain, Chain: &ChainEvent{Verb: verb, ID: id, Chain: chain}})
}

// PublishTrigger is a convenience for redacted external-trigger events.
func (b *Broker) PublishTrigger(verb Verb, id string, trigger *domain.ExternalTrigger) {
	if trigger != nil {
		copy := *trigger
		copy.Key = ""
		trigger = &copy
	}
	b.Publish(Event{Kind: KindTrigger, Trigger: &TriggerEvent{Verb: verb, ID: id, Trigger: trigger}})
}

// PublishTriggerSet is a convenience for redacted Trigger Set events.
func (b *Broker) PublishTriggerSet(verb Verb, id string, set *domain.TriggerSet) {
	if set != nil {
		copy := *set
		copy.Members = append([]domain.ExternalTrigger(nil), set.Members...)
		for index := range copy.Members {
			copy.Members[index].Key = ""
		}
		set = &copy
	}
	b.Publish(Event{Kind: KindTriggerSet, TriggerSet: &TriggerSetEvent{Verb: verb, ID: id, TriggerSet: set}})
}

// PublishWatcher is a convenience for path-free watcher lifecycle and health events.
func (b *Broker) PublishWatcher(verb Verb, id, name string, health *domain.WatcherHealth) {
	b.Publish(Event{Kind: KindWatcher, Watcher: &WatcherEvent{Verb: verb, ID: id, Name: name, Health: health}})
}

// SubscriberCount reports the number of active subscribers (for tests/metrics).
func (b *Broker) SubscriberCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subs)
}
