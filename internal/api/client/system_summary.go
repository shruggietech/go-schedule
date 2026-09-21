package client

import (
	"context"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// SystemSummary returns one bounded, read-only operational observation from the selected daemon.
func (c *Client) SystemSummary(ctx context.Context) (domain.SystemSummary, error) {
	var result domain.SystemSummary
	if err := c.get(ctx, "/v1/system-summary", &result); err != nil {
		return domain.SystemSummary{}, err
	}
	return result, nil
}
