// Package bundles exposes the desktop workflow for portable automation bundles.
package bundles

import "github.com/shruggietech/go-schedule/internal/bundle"

type Result struct {
	Action      string         `json:"action"`
	Outcome     string         `json:"outcome"`
	Message     string         `json:"message"`
	Document    string         `json:"document,omitempty"`
	Plan        *bundle.Plan   `json:"plan,omitempty"`
	CompareOnly bool           `json:"compare_only,omitempty"`
	Issues      []bundle.Issue `json:"issues,omitempty"`
	Items       []bundle.Item  `json:"items,omitempty"`
}
