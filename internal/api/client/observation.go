package client

import (
	"context"
	"errors"
	"strings"
	"time"
)

// RequestObservation contains only authored route and outcome categories, never request content.
type RequestObservation struct {
	Method     string `json:"method"`
	Route      string `json:"route"`
	DurationMS int64  `json:"duration_ms"`
	Outcome    string `json:"outcome"`
}
type requestObserverKey struct{}

// WithRequestObserver scopes diagnostics to this context. The callback must return promptly.
func WithRequestObserver(ctx context.Context, observer func(RequestObservation)) context.Context {
	return context.WithValue(ctx, requestObserverKey{}, observer)
}

func observeRequest(ctx context.Context, method, path string) func(error) {
	observer, _ := ctx.Value(requestObserverKey{}).(func(RequestObservation))
	if observer == nil {
		return func(error) {}
	}
	start := time.Now()
	route := observationRoute(path)
	return func(err error) {
		outcome := "completed"
		if err != nil {
			outcome = "failed"
			var status *StatusError
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				outcome = "deadline_exceeded"
			case errors.Is(err, context.Canceled):
				outcome = "canceled"
			case errors.As(err, &status):
				outcome = "api_rejected"
			}
		}
		observer(RequestObservation{method, route, time.Since(start).Milliseconds(), outcome})
	}
}

// Only fixed API vocabulary is emitted. IDs, query values and arbitrary path text are excluded.
func observationRoute(path string) string {
	path, _, _ = strings.Cut(path, "?")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 8 {
		return "/{unrecognized}"
	}
	for i, part := range parts {
		if !observationVocabulary[part] {
			parts[i] = "{id}"
		}
	}
	return "/" + strings.Join(parts, "/")
}

var observationVocabulary = func() map[string]bool {
	vocabulary := map[string]bool{}
	for _, part := range strings.Fields("v1 health manifest runtime-info tasks groups chains triggers trigger-sets filesystem-watchers notification notification-channels notification-deliveries channels assignments deliveries test enabled authorization access current actors credentials pairings mcp http sessions logs runs alerts active settings schedules preview run-now enable disable rotate revoke reveal fire retarget ack system summary search bundles export import plan audit") {
		vocabulary[part] = true
	}
	return vocabulary
}()
