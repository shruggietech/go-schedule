package automation

import "context"

type joinSupportKey struct{}

// WithJoinSupport captures the selected scheduler's advertised join capability for mutation refreshes.
func WithJoinSupport(ctx context.Context, supported bool) context.Context {
	return context.WithValue(ctx, joinSupportKey{}, supported)
}

func joinSupport(ctx context.Context) bool {
	supported, selected := ctx.Value(joinSupportKey{}).(bool)
	// Direct service users retain the current-daemon contract; the Wails facade always captures its selected manifest.
	return !selected || supported
}
