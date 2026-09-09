package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func (c *Client) ListActors(ctx context.Context) ([]domain.Actor, error) {
	var out struct {
		Actors []domain.Actor `json:"actors"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/access/actors", nil, &out)
	return out.Actors, err
}

func (c *Client) CreateActor(ctx context.Context, request server.ActorCreateRequest) (domain.Actor, error) {
	var out domain.Actor
	err := c.do(ctx, http.MethodPost, "/v1/access/actors", request, &out)
	return out, err
}

func (c *Client) UpdateActor(ctx context.Context, id string, request server.ActorUpdateRequest) (domain.Actor, error) {
	var out domain.Actor
	err := c.do(ctx, http.MethodPatch, "/v1/access/actors/"+url.PathEscape(id), request, &out)
	return out, err
}

func (c *Client) RevokeActor(ctx context.Context, id string) (domain.Actor, error) {
	var out domain.Actor
	err := c.do(ctx, http.MethodPost, "/v1/access/actors/"+url.PathEscape(id)+"/revoke", nil, &out)
	return out, err
}

func auditQueryValues(query domain.AuditQuery) url.Values {
	values := url.Values{}
	if query.ActorID != "" {
		values.Set("actor_id", query.ActorID)
	}
	if query.Operation != "" {
		values.Set("operation", query.Operation)
	}
	if query.Result != "" {
		values.Set("result", string(query.Result))
	}
	if query.Since != nil {
		values.Set("since", query.Since.Format(timeFormat))
	}
	if query.Until != nil {
		values.Set("until", query.Until.Format(timeFormat))
	}
	if query.Limit > 0 {
		values.Set("limit", strconv.Itoa(query.Limit))
	}
	return values
}

const timeFormat = "2006-01-02T15:04:05.999999999Z07:00"

func (c *Client) ListAudit(ctx context.Context, query domain.AuditQuery) ([]domain.AuditEvent, error) {
	var out struct {
		Events []domain.AuditEvent `json:"events"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/audit", auditQueryValues(query)), nil, &out)
	return out.Events, err
}

func (c *Client) ExportAudit(ctx context.Context, query domain.AuditQuery) ([]byte, error) {
	path := withQuery("/v1/audit/export", auditQueryValues(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, NewConnectionError("GET "+path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, decodeStatusError(resp, http.MethodGet, path)
	}
	return io.ReadAll(resp.Body)
}

func decodeStatusError(resp *http.Response, method, path string) error {
	var apiErr server.APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error.Message != "" {
		return &StatusError{Code: apiErr.Error.Code, Field: apiErr.Error.Field, Message: apiErr.Error.Message}
	}
	return &StatusError{Code: server.CodeInternal, Message: method + " " + path + ": unexpected status"}
}
