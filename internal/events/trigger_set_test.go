package events

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestPublishTriggerSetEmitsOneRedactedEvent(t *testing.T) {
	broker := NewBroker()
	stream, cancel := broker.Subscribe()
	defer cancel()
	set := domain.TriggerSet{ID: "set", Name: "Callers", TargetTaskID: "task", Members: []domain.ExternalTrigger{{ID: "trigger", Key: "gst_secret", SetID: "set", SetPosition: 1}}}
	broker.PublishTriggerSet(VerbUpdated, set.ID, &set)
	select {
	case event := <-stream:
		encoded, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		if event.Kind != KindTriggerSet || event.TriggerSet == nil || event.TriggerSet.TriggerSet == nil || strings.Contains(string(encoded), "gst_secret") {
			t.Fatalf("event=%s", encoded)
		}
		select {
		case extra := <-stream:
			t.Fatalf("unexpected second event=%+v", extra)
		default:
		}
	case <-time.After(time.Second):
		t.Fatal("Trigger Set event was not delivered")
	}
}
