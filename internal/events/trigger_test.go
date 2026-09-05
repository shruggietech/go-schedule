package events

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestPublishTriggerRedactsKey(t *testing.T) {
	broker := NewBroker()
	events, cancel := broker.Subscribe()
	defer cancel()
	trigger := domain.ExternalTrigger{ID: "id", Name: "hook", Key: "gst_secret", TargetTaskID: "task", Enabled: true}
	broker.PublishTrigger(VerbCreated, trigger.ID, &trigger)
	select {
	case event := <-events:
		if event.Kind != KindTrigger || event.Trigger == nil || event.Trigger.Trigger == nil {
			t.Fatalf("event = %+v", event)
		}
		if event.Trigger.Trigger.Key != "" {
			t.Fatal("trigger event disclosed key")
		}
	case <-time.After(time.Second):
		t.Fatal("trigger event was not delivered")
	}
}
