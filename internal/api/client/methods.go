package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shruggietech/go-scheduler/internal/api/server"
	"github.com/shruggietech/go-scheduler/internal/domain"
)

// CreateTask creates a task and returns its detail.
func (c *Client) CreateTask(ctx context.Context, req server.TaskCreateRequest) (server.TaskResponse, error) {
	var out server.TaskResponse
	err := c.do(ctx, http.MethodPost, "/v1/tasks", req, &out)
	return out, err
}

// ListTasks lists tasks, optionally filtered by group and state.
func (c *Client) ListTasks(ctx context.Context, group, state string) ([]domain.Task, error) {
	q := url.Values{}
	if group != "" {
		q.Set("group", group)
	}
	if state != "" {
		q.Set("state", state)
	}
	var out struct {
		Tasks []domain.Task `json:"tasks"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/tasks", q), nil, &out)
	return out.Tasks, err
}

// GetTask returns a task's detail.
func (c *Client) GetTask(ctx context.Context, id string) (server.TaskResponse, error) {
	var out server.TaskResponse
	err := c.do(ctx, http.MethodGet, "/v1/tasks/"+id, nil, &out)
	return out, err
}

// UpdateTask applies partial changes to a task.
func (c *Client) UpdateTask(ctx context.Context, id string, req server.TaskUpdateRequest) (server.TaskResponse, error) {
	var out server.TaskResponse
	err := c.do(ctx, http.MethodPatch, "/v1/tasks/"+id, req, &out)
	return out, err
}

// DeleteTask deletes a task.
func (c *Client) DeleteTask(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/tasks/"+id, nil, nil)
}

// SetTaskEnabled enables or disables a task.
func (c *Client) SetTaskEnabled(ctx context.Context, id string, enabled bool) error {
	action := "disable"
	if enabled {
		action = "enable"
	}
	return c.do(ctx, http.MethodPost, "/v1/tasks/"+id+"/"+action, nil, nil)
}

// RunNow triggers an immediate run.
func (c *Client) RunNow(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/v1/tasks/"+id+"/run-now", nil, nil)
}

// Preview returns the RRULE, summary, and next runs for a schedule expression.
func (c *Client) Preview(ctx context.Context, req server.PreviewRequest) (server.PreviewResponse, error) {
	var out server.PreviewResponse
	err := c.do(ctx, http.MethodPost, "/v1/schedules/preview", req, &out)
	return out, err
}

// ListRuns returns run history (optionally for one task).
func (c *Client) ListRuns(ctx context.Context, taskID string, limit int) ([]domain.Run, error) {
	q := url.Values{}
	if taskID != "" {
		q.Set("task", taskID)
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var out struct {
		Runs []domain.Run `json:"runs"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/runs", q), nil, &out)
	return out.Runs, err
}

// ListAlerts returns alerts, optionally only unacknowledged.
func (c *Client) ListAlerts(ctx context.Context, unacked bool) ([]domain.Alert, error) {
	q := url.Values{}
	if unacked {
		q.Set("unacked", "true")
	}
	var out struct {
		Alerts []domain.Alert `json:"alerts"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/alerts", q), nil, &out)
	return out.Alerts, err
}

// AckAlert acknowledges an alert.
func (c *Client) AckAlert(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodPost, "/v1/alerts/"+id+"/ack", nil, nil)
}

func withQuery(path string, q url.Values) string {
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}

// do performs a request with an optional JSON body and decodes an optional JSON
// response, surfacing the API error envelope.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var rdr *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("api: %s %s: %w (is the daemon running?)", method, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var apiErr server.APIError
		if decErr := json.NewDecoder(resp.Body).Decode(&apiErr); decErr == nil && apiErr.Error.Message != "" {
			return &StatusError{Code: apiErr.Error.Code, Field: apiErr.Error.Field, Message: apiErr.Error.Message}
		}
		return &StatusError{Code: server.CodeInternal, Message: fmt.Sprintf("%s %s: status %d", method, path, resp.StatusCode)}
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("api: decode %s: %w", path, err)
		}
	}
	return nil
}
