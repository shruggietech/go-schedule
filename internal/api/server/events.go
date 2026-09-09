package server

import (
	"encoding/json"
	"net/http"

	"github.com/shruggietech/go-schedule/internal/events"
)

// handleEvents streams run-state changes and new alerts to the client as
// Server-Sent Events, so the GUI can surface updates within seconds without
// polling. Each event is a JSON object on a single `data:` line.
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	if s.broker == nil {
		writeError(w, http.StatusServiceUnavailable, CodeInternal, "", "event stream unavailable")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, CodeInternal, "", "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch, cancel := s.broker.Subscribe()
	defer cancel()

	// Initial comment so clients know the stream is open.
	_, _ = w.Write([]byte(": connected\n\n"))
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			payload := any(ev)
			if r.URL.Query().Get("observation") == "true" {
				projected, visible := remoteEventProjection(ev)
				if !visible {
					continue
				}
				payload = projected
			}
			data, err := json.Marshal(payload)
			if err != nil {
				continue
			}
			if _, err := w.Write([]byte("event: " + string(ev.Kind) + "\ndata: ")); err != nil {
				return
			}
			if _, err := w.Write(data); err != nil {
				return
			}
			if _, err := w.Write([]byte("\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

type remoteEvent struct {
	Kind       events.Kind `json:"kind"`
	ResourceID string      `json:"resource_id"`
	Verb       events.Verb `json:"verb,omitempty"`
}

func remoteEventProjection(event events.Event) (remoteEvent, bool) {
	projected := remoteEvent{Kind: event.Kind}
	switch {
	case event.Run != nil:
		projected.ResourceID = event.Run.ID
	case event.Alert != nil:
		projected.ResourceID = event.Alert.ID
	case event.Task != nil:
		projected.ResourceID, projected.Verb = event.Task.ID, event.Task.Verb
	case event.Group != nil:
		projected.ResourceID, projected.Verb = event.Group.ID, event.Group.Verb
	case event.Chain != nil:
		projected.ResourceID, projected.Verb = event.Chain.ID, event.Chain.Verb
	case event.Trigger != nil:
		projected.ResourceID, projected.Verb = event.Trigger.ID, event.Trigger.Verb
	case event.TriggerSet != nil:
		projected.ResourceID, projected.Verb = event.TriggerSet.ID, event.TriggerSet.Verb
	case event.Watcher != nil:
		projected.ResourceID, projected.Verb = event.Watcher.ID, event.Watcher.Verb
	default:
		return remoteEvent{}, false
	}
	return projected, projected.ResourceID != ""
}
