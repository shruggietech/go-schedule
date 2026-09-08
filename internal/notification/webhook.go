// Package notification validates webhook destinations and delivers durable
// notification work without participating in scheduler task execution.
package notification

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/buildinfo"
	"github.com/shruggietech/go-schedule/internal/domain"
)

// RequestTimeout bounds one outbound webhook attempt.
const RequestTimeout = 5 * time.Second

// ValidateEndpoint validates a webhook URL and returns its safe scheme-and-host summary.
func ValidateEndpoint(raw string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("endpoint must be an absolute HTTP URL")
	}
	if parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return "", fmt.Errorf("endpoint must use https with a host")
	}
	if parsed.Fragment != "" {
		return "", fmt.Errorf("endpoint must not contain a fragment")
	}
	if parsed.Scheme == "http" && !isLoopbackHost(parsed.Hostname()) {
		return "", fmt.Errorf("endpoint must use https unless its host is loopback")
	}
	return parsed.Scheme + "://" + parsed.Host, nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// Sender performs one webhook attempt. It never follows redirects.
type Sender struct {
	client *http.Client
}

// NewSender constructs a webhook sender using transport or the default transport.
func NewSender(transport http.RoundTripper) *Sender {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &Sender{client: &http.Client{
		Transport: transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}}
}

// Send delivers one attempt and returns its HTTP status and safe diagnostic.
func (s *Sender) Send(ctx context.Context, delivery domain.NotificationDelivery) (int, string, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(attemptCtx, http.MethodPost, delivery.Endpoint, bytes.NewReader(delivery.Payload))
	if err != nil {
		return 0, "request construction failed", fmt.Errorf("construct webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-schedule/"+buildinfo.Version)
	req.Header.Set("X-Go-Schedule-Delivery", delivery.ID)
	req.Header.Set("X-Go-Schedule-Event", string(delivery.EventKind))
	if delivery.Authorization != "" {
		req.Header.Set("Authorization", delivery.Authorization)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		if errorsIsTimeout(attemptCtx, err) {
			return 0, "request timed out", err
		}
		return 0, "transport failed", err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		diagnostic := fmt.Sprintf("receiver returned HTTP %d", resp.StatusCode)
		return resp.StatusCode, diagnostic, fmt.Errorf("%s", diagnostic)
	}
	return resp.StatusCode, "", nil
}

func errorsIsTimeout(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return true
	}
	type timeout interface{ Timeout() bool }
	value, ok := err.(timeout)
	return ok && value.Timeout()
}
