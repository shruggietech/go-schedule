package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
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

// CreateChain creates a task-completion relationship.
func (c *Client) CreateChain(ctx context.Context, req server.ChainCreateRequest) (domain.CompletionChain, error) {
	var out domain.CompletionChain
	err := c.do(ctx, http.MethodPost, "/v1/chains", req, &out)
	return out, err
}

// ListChains returns all task-completion relationships.
func (c *Client) ListChains(ctx context.Context) ([]domain.CompletionChain, error) {
	var out struct {
		Chains []domain.CompletionChain `json:"chains"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/chains", nil, &out)
	return out.Chains, err
}

// GetChain returns one task-completion relationship.
func (c *Client) GetChain(ctx context.Context, id string) (domain.CompletionChain, error) {
	var out domain.CompletionChain
	err := c.do(ctx, http.MethodGet, "/v1/chains/"+id, nil, &out)
	return out, err
}

// UpdateChain changes one or more fields of a task-completion relationship.
func (c *Client) UpdateChain(ctx context.Context, id string, req server.ChainUpdateRequest) (domain.CompletionChain, error) {
	var out domain.CompletionChain
	err := c.do(ctx, http.MethodPatch, "/v1/chains/"+id, req, &out)
	return out, err
}

// DeleteChain removes a task-completion relationship.
func (c *Client) DeleteChain(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/chains/"+id, nil, nil)
}

func (c *Client) CreateTrigger(ctx context.Context, req server.TriggerCreateRequest) (server.TriggerSecretResponse, error) {
	var out server.TriggerSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/triggers", req, &out)
	return out, err
}

func (c *Client) ListTriggers(ctx context.Context) ([]server.TriggerResponse, error) {
	var out struct {
		Triggers []server.TriggerResponse `json:"triggers"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/triggers", nil, &out)
	return out.Triggers, err
}

func (c *Client) GetTrigger(ctx context.Context, id string) (server.TriggerResponse, error) {
	var out server.TriggerResponse
	err := c.do(ctx, http.MethodGet, "/v1/triggers/"+url.PathEscape(id), nil, &out)
	return out, err
}

func (c *Client) UpdateTrigger(ctx context.Context, id string, req server.TriggerUpdateRequest) (server.TriggerResponse, error) {
	var out server.TriggerResponse
	err := c.do(ctx, http.MethodPatch, "/v1/triggers/"+url.PathEscape(id), req, &out)
	return out, err
}

func (c *Client) DeleteTrigger(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/triggers/"+url.PathEscape(id), nil, nil)
}

func (c *Client) SetTriggerEnabled(ctx context.Context, id string, enabled bool) (server.TriggerResponse, error) {
	action := "disable"
	if enabled {
		action = "enable"
	}
	var out server.TriggerResponse
	err := c.do(ctx, http.MethodPost, "/v1/triggers/"+url.PathEscape(id)+"/"+action, nil, &out)
	return out, err
}

func (c *Client) RotateTrigger(ctx context.Context, id string) (server.TriggerSecretResponse, error) {
	var out server.TriggerSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/triggers/"+url.PathEscape(id)+"/rotate", nil, &out)
	return out, err
}

func (c *Client) RevealTrigger(ctx context.Context, id string) (server.TriggerSecretResponse, error) {
	var out server.TriggerSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/triggers/"+url.PathEscape(id)+"/reveal", nil, &out)
	return out, err
}

func (c *Client) FireTrigger(ctx context.Context, key string) error {
	return c.do(ctx, http.MethodPost, "/v1/triggers/fire", map[string]string{"key": key}, nil)
}

func (c *Client) CreateFilesystemWatcher(ctx context.Context, req server.FilesystemWatcherCreateRequest) (server.FilesystemWatcherResponse, error) {
	var out server.FilesystemWatcherResponse
	err := c.do(ctx, http.MethodPost, "/v1/filesystem-watchers", req, &out)
	return out, err
}

func (c *Client) ListFilesystemWatchers(ctx context.Context) ([]server.FilesystemWatcherResponse, error) {
	var out struct {
		FilesystemWatchers []server.FilesystemWatcherResponse `json:"filesystem_watchers"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/filesystem-watchers", nil, &out)
	return out.FilesystemWatchers, err
}

func (c *Client) GetFilesystemWatcher(ctx context.Context, id string) (server.FilesystemWatcherResponse, error) {
	var out server.FilesystemWatcherResponse
	err := c.do(ctx, http.MethodGet, "/v1/filesystem-watchers/"+url.PathEscape(id), nil, &out)
	return out, err
}

func (c *Client) UpdateFilesystemWatcher(ctx context.Context, id string, req server.FilesystemWatcherUpdateRequest) (server.FilesystemWatcherResponse, error) {
	var out server.FilesystemWatcherResponse
	err := c.do(ctx, http.MethodPatch, "/v1/filesystem-watchers/"+url.PathEscape(id), req, &out)
	return out, err
}

func (c *Client) SetFilesystemWatcherEnabled(ctx context.Context, id string, enabled bool) (server.FilesystemWatcherResponse, error) {
	action := "disable"
	if enabled {
		action = "enable"
	}
	var out server.FilesystemWatcherResponse
	err := c.do(ctx, http.MethodPost, "/v1/filesystem-watchers/"+url.PathEscape(id)+"/"+action, nil, &out)
	return out, err
}

func (c *Client) DeleteFilesystemWatcher(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/filesystem-watchers/"+url.PathEscape(id), nil, nil)
}

// CreateTriggerSet atomically creates a set and returns its ordered secrets.
func (c *Client) CreateTriggerSet(ctx context.Context, req server.TriggerSetCreateRequest) (server.TriggerSetSecretResponse, error) {
	var out server.TriggerSetSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/trigger-sets", req, &out)
	return out, err
}

// ListTriggerSets returns every Trigger Set without raw keys.
func (c *Client) ListTriggerSets(ctx context.Context) ([]server.TriggerSetResponse, error) {
	var out struct {
		TriggerSets []server.TriggerSetResponse `json:"trigger_sets"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/trigger-sets", nil, &out)
	return out.TriggerSets, err
}

// GetTriggerSet returns one Trigger Set without raw keys.
func (c *Client) GetTriggerSet(ctx context.Context, id string) (server.TriggerSetResponse, error) {
	var out server.TriggerSetResponse
	err := c.do(ctx, http.MethodGet, "/v1/trigger-sets/"+url.PathEscape(id), nil, &out)
	return out, err
}

// RevealTriggerSet returns the current ordered member secrets explicitly.
func (c *Client) RevealTriggerSet(ctx context.Context, id string) (server.TriggerSetSecretResponse, error) {
	var out server.TriggerSetSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/trigger-sets/"+url.PathEscape(id)+"/reveal", nil, &out)
	return out, err
}

// RetargetTriggerSet atomically changes every member target.
func (c *Client) RetargetTriggerSet(ctx context.Context, id, taskID string) (server.TriggerSetResponse, error) {
	var out server.TriggerSetResponse
	err := c.do(ctx, http.MethodPatch, "/v1/trigger-sets/"+url.PathEscape(id), server.TriggerSetRetargetRequest{TargetTaskID: taskID}, &out)
	return out, err
}

// SetTriggerSetEnabled atomically changes every member's enabled state.
func (c *Client) SetTriggerSetEnabled(ctx context.Context, id string, enabled bool) (server.TriggerSetResponse, error) {
	action := "disable"
	if enabled {
		action = "enable"
	}
	var out server.TriggerSetResponse
	err := c.do(ctx, http.MethodPost, "/v1/trigger-sets/"+url.PathEscape(id)+"/"+action, nil, &out)
	return out, err
}

// RotateTriggerSet atomically replaces every member key and returns replacements.
func (c *Client) RotateTriggerSet(ctx context.Context, id string) (server.TriggerSetSecretResponse, error) {
	var out server.TriggerSetSecretResponse
	err := c.do(ctx, http.MethodPost, "/v1/trigger-sets/"+url.PathEscape(id)+"/rotate", nil, &out)
	return out, err
}

// DeleteTriggerSet atomically removes a set and every member.
func (c *Client) DeleteTriggerSet(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/trigger-sets/"+url.PathEscape(id), nil, nil)
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

// GetRun returns one run by exact identity.
func (c *Client) GetRun(ctx context.Context, id string) (domain.Run, error) {
	var out domain.Run
	err := c.do(ctx, http.MethodGet, "/v1/runs/"+url.PathEscape(id), nil, &out)
	return out, err
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
		return NewConnectionError(method+" "+path, err)
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
