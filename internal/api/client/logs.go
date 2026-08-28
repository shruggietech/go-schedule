package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// ListLogs returns recent log records and daemon log metadata. Records are
// newest first; severity ("" | info | warning | error) and limit (0 = server
// default) filter the results.
func (c *Client) ListLogs(ctx context.Context, severity string, limit int) (server.LogsResponse, error) {
	q := url.Values{}
	if severity != "" {
		q.Set("severity", severity)
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var out server.LogsResponse
	err := c.do(ctx, http.MethodGet, withQuery("/v1/logs", q), nil, &out)
	return out, err
}
