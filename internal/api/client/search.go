package client

import (
	"context"
	"net/url"
	"strconv"

	"github.com/shruggietech/go-schedule/internal/domain"
)

// Search returns a bounded, secret-free search observation from the selected daemon.
func (c *Client) Search(ctx context.Context, query string, kinds []domain.SearchKind, limit int) (domain.DaemonSearch, error) {
	values := url.Values{"q": []string{query}}
	for _, kind := range kinds {
		values.Add("kind", string(kind))
	}
	if limit > 0 {
		values.Set("limit", strconv.Itoa(limit))
	}
	var result domain.DaemonSearch
	if err := c.get(ctx, "/v1/search?"+values.Encode(), &result); err != nil {
		return domain.DaemonSearch{}, err
	}
	return result, nil
}
