package client

import (
	"context"
	"net/http"

	"github.com/shruggietech/go-scheduler/internal/api/server"
	"github.com/shruggietech/go-scheduler/internal/domain"
)

// CreateTrigger creates an event trigger.
func (c *Client) CreateTrigger(ctx context.Context, req server.TriggerCreateRequest) (domain.Trigger, error) {
	var tr domain.Trigger
	err := c.do(ctx, http.MethodPost, "/v1/triggers", req, &tr)
	return tr, err
}

// ListTriggers lists all triggers.
func (c *Client) ListTriggers(ctx context.Context) ([]domain.Trigger, error) {
	var out struct {
		Triggers []domain.Trigger `json:"triggers"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/triggers", nil, &out)
	return out.Triggers, err
}

// DeleteTrigger removes a trigger.
func (c *Client) DeleteTrigger(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/triggers/"+id, nil, nil)
}
