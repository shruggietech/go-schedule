package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// CreateNotificationChannel creates one write-only webhook destination.
func (c *Client) CreateNotificationChannel(ctx context.Context, req server.NotificationChannelCreateRequest) (domain.NotificationChannel, error) {
	var out domain.NotificationChannel
	err := c.do(ctx, http.MethodPost, "/v1/notification-channels", req, &out)
	return out, err
}

// ListNotificationChannels returns redacted channel metadata.
func (c *Client) ListNotificationChannels(ctx context.Context) ([]domain.NotificationChannel, error) {
	var out struct {
		Channels []domain.NotificationChannel `json:"notification_channels"`
	}
	err := c.do(ctx, http.MethodGet, "/v1/notification-channels", nil, &out)
	return out.Channels, err
}

// GetNotificationChannel returns one redacted channel.
func (c *Client) GetNotificationChannel(ctx context.Context, id string) (domain.NotificationChannel, error) {
	var out domain.NotificationChannel
	err := c.do(ctx, http.MethodGet, "/v1/notification-channels/"+url.PathEscape(id), nil, &out)
	return out, err
}

// UpdateNotificationChannel updates mutable channel metadata.
func (c *Client) UpdateNotificationChannel(ctx context.Context, id string, req server.NotificationChannelUpdateRequest) (domain.NotificationChannel, error) {
	var out domain.NotificationChannel
	err := c.do(ctx, http.MethodPatch, "/v1/notification-channels/"+url.PathEscape(id), req, &out)
	return out, err
}

// SetNotificationChannelEnabled enables or disables new delivery creation.
func (c *Client) SetNotificationChannelEnabled(ctx context.Context, id string, enabled bool) (domain.NotificationChannel, error) {
	action := "disable"
	if enabled {
		action = "enable"
	}
	var out domain.NotificationChannel
	err := c.do(ctx, http.MethodPost, "/v1/notification-channels/"+url.PathEscape(id)+"/"+action, nil, &out)
	return out, err
}

// RotateNotificationChannelAuthorization replaces or clears authorization.
func (c *Client) RotateNotificationChannelAuthorization(ctx context.Context, id, authorization string) (domain.NotificationChannel, error) {
	var out domain.NotificationChannel
	err := c.do(ctx, http.MethodPost, "/v1/notification-channels/"+url.PathEscape(id)+"/rotate", server.NotificationChannelRotateRequest{Authorization: authorization}, &out)
	return out, err
}

// TestNotificationChannel queues a transport-equivalent test event.
func (c *Client) TestNotificationChannel(ctx context.Context, id string) (domain.NotificationDelivery, error) {
	var out domain.NotificationDelivery
	err := c.do(ctx, http.MethodPost, "/v1/notification-channels/"+url.PathEscape(id)+"/test", nil, &out)
	return out, err
}

// DeleteNotificationChannel removes a channel and unfinished work.
func (c *Client) DeleteNotificationChannel(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/notification-channels/"+url.PathEscape(id), nil, nil)
}

// ReplaceTaskNotificationAssignments replaces a task's direct policy.
func (c *Client) ReplaceTaskNotificationAssignments(ctx context.Context, id string, req server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	return c.replaceNotificationAssignments(ctx, "/v1/tasks/"+url.PathEscape(id)+"/notifications", req)
}

// ReplaceGroupNotificationAssignments replaces a group's direct policy.
func (c *Client) ReplaceGroupNotificationAssignments(ctx context.Context, id string, req server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	return c.replaceNotificationAssignments(ctx, "/v1/groups/"+url.PathEscape(id)+"/notifications", req)
}

func (c *Client) replaceNotificationAssignments(ctx context.Context, path string, req server.NotificationAssignmentsRequest) ([]domain.NotificationAssignment, error) {
	var out struct {
		Assignments []domain.NotificationAssignment `json:"assignments"`
	}
	err := c.do(ctx, http.MethodPut, path, req, &out)
	return out.Assignments, err
}

// ListTaskNotificationAssignments returns a task's directly configured policy.
func (c *Client) ListTaskNotificationAssignments(ctx context.Context, id string) ([]domain.NotificationAssignment, error) {
	return c.listNotificationAssignments(ctx, "/v1/tasks/"+url.PathEscape(id)+"/notifications")
}

// ListGroupNotificationAssignments returns a group's directly configured policy.
func (c *Client) ListGroupNotificationAssignments(ctx context.Context, id string) ([]domain.NotificationAssignment, error) {
	return c.listNotificationAssignments(ctx, "/v1/groups/"+url.PathEscape(id)+"/notifications")
}

func (c *Client) listNotificationAssignments(ctx context.Context, path string) ([]domain.NotificationAssignment, error) {
	var out struct {
		Assignments []domain.NotificationAssignment `json:"assignments"`
	}
	err := c.do(ctx, http.MethodGet, path, nil, &out)
	return out.Assignments, err
}

// EffectiveTaskNotificationPolicy returns the selected source and assignments.
func (c *Client) EffectiveTaskNotificationPolicy(ctx context.Context, id string) (domain.EffectiveNotificationPolicy, error) {
	var out domain.EffectiveNotificationPolicy
	err := c.do(ctx, http.MethodGet, "/v1/tasks/"+url.PathEscape(id)+"/notifications/effective", nil, &out)
	return out, err
}

// ListNotificationDeliveries returns filtered redacted evidence.
func (c *Client) ListNotificationDeliveries(ctx context.Context, filter domain.NotificationDeliveryFilter) ([]domain.NotificationDelivery, error) {
	q := url.Values{}
	if filter.ChannelID != "" {
		q.Set("channel", filter.ChannelID)
	}
	if filter.TaskID != "" {
		q.Set("task", filter.TaskID)
	}
	if filter.RunID != "" {
		q.Set("run", filter.RunID)
	}
	if filter.State != "" {
		q.Set("state", string(filter.State))
	}
	if filter.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", filter.Limit))
	}
	var out struct {
		Deliveries []domain.NotificationDelivery `json:"notification_deliveries"`
	}
	err := c.do(ctx, http.MethodGet, withQuery("/v1/notification-deliveries", q), nil, &out)
	return out.Deliveries, err
}
