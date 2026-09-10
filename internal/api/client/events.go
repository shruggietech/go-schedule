package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/events"
)

// RemoteEvent is the secret-free event identity exposed by the HTTPS API.
type RemoteEvent struct {
	Kind       string `json:"kind"`
	ResourceID string `json:"resource_id"`
	Verb       string `json:"verb,omitempty"`
}

// GetCalendar returns calendar occurrences in [from, to].
func (c *Client) GetCalendar(ctx context.Context, from, to time.Time) (server.CalendarResponse, error) {
	var out server.CalendarResponse
	path := "/v1/calendar?from=" + from.UTC().Format(time.RFC3339) + "&to=" + to.UTC().Format(time.RFC3339)
	err := c.do(ctx, http.MethodGet, path, nil, &out)
	return out, err
}

// StreamEvents opens the SSE event stream and invokes onEvent for each event
// until ctx is cancelled or the stream ends. It blocks; run it in a goroutine.
func (c *Client) StreamEvents(ctx context.Context, onEvent func(events.Event)) error {
	target := c.target()
	req, err := target.newRequest(ctx, http.MethodGet, "/v1/events", nil)
	if err != nil {
		return err
	}
	resp, err := target.http.Do(req)
	if err != nil {
		return NewConnectionError("GET /v1/events", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		var apiErr server.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error.Message != "" {
			return &StatusError{Code: apiErr.Error.Code, Field: apiErr.Error.Field, Message: apiErr.Error.Message}
		}
		return &StatusError{Code: server.CodeInternal, Message: fmt.Sprintf("GET /v1/events: status %d", resp.StatusCode)}
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return NewConnectionError("GET /v1/events", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var ev events.Event
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &ev); err == nil {
			onEvent(ev)
		}
	}
}

// StreamRemoteEvents consumes the secret-free remote event projection.
func (c *Client) StreamRemoteEvents(ctx context.Context, onEvent func(RemoteEvent)) error {
	target := c.target()
	req, err := target.newRequest(ctx, http.MethodGet, "/v1/events", nil)
	if err != nil {
		return err
	}
	resp, err := target.http.Do(req)
	if err != nil {
		return NewConnectionError("GET /v1/events", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		return decodeStatusError(resp, http.MethodGet, "/v1/events")
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			return NewConnectionError("GET /v1/events", readErr)
		}
		line = strings.TrimRight(line, "\r\n")
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var event RemoteEvent
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event) == nil && event.Kind != "" && event.ResourceID != "" {
			onEvent(event)
		}
	}
}
