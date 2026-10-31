package clievaluate

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/notification"
)

type ownedRoundTrip func(*http.Request) (*http.Response, error)

func (f ownedRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func notificationTransport(ctx context.Context) execution {
	delivery := domain.NotificationDelivery{ID: "owned-delivery", Endpoint: "http://127.0.0.1/owned", Payload: []byte(`{}`)}
	observations := []NamedValue{}
	good := true
	for _, status := range []int{http.StatusServiceUnavailable, http.StatusTemporaryRedirect} {
		calls := 0
		sender := notification.NewSender(ownedRoundTrip(func(r *http.Request) (*http.Response, error) {
			calls++
			return &http.Response{StatusCode: status, Header: http.Header{"Location": []string{"http://127.0.0.1/forbidden-redirect"}}, Body: io.NopCloser(strings.NewReader("owned rejected response")), Request: r}, nil
		}))
		actual, diagnostic, err := sender.Send(ctx, delivery)
		good = good && err != nil && actual == status && calls == 1 && strings.Contains(diagnostic, "HTTP")
		observations = append(observations, NamedValue{Name: http.StatusText(status), Value: diagnostic})
	}
	timeout, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	sender := notification.NewSender(ownedRoundTrip(func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() }))
	status, diagnostic, err := sender.Send(timeout, delivery)
	good = good && err != nil && status == 0 && diagnostic == "request timed out"
	observations = append(observations, NamedValue{Name: "bounded timeout", Value: diagnostic})
	result := policyResult(observations, good, nil)
	result.Logs = mustObservation(observations)
	return result
}
