package events

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestWatcherEventOmitsFilesystemPath(t *testing.T) {
	broker := NewBroker()
	events, cancel := broker.Subscribe()
	defer cancel()
	health := domain.WatcherHealth{State: domain.WatcherDegraded, Reason: "configured observation root is missing"}
	broker.PublishWatcher(VerbUpdated, "watcher-1", "incoming", &health)
	payload, err := json.Marshal(<-events)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), `C:\\sensitive`) || strings.Contains(string(payload), `"path"`) {
		t.Fatalf("event disclosed path: %s", payload)
	}
}
