package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestNotificationClientUsesVersionedPathsAndFilters(t *testing.T) {
	var method, path, query, body string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		method, path, query = request.Method, request.URL.Path, request.URL.RawQuery
		data, _ := io.ReadAll(request.Body)
		body = string(data)
		response := `{"assignments":[]}`
		if path == "/v1/notification-deliveries" {
			response = `{"notification_deliveries":[]}`
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}, nil
	})}}
	_, err := client.ReplaceTaskNotificationAssignments(context.Background(), "task one", server.NotificationAssignmentsRequest{Assignments: []server.NotificationAssignmentInput{{ChannelID: "channel", OnFailure: true}}})
	if err != nil || method != http.MethodPut || path != "/v1/tasks/task one/notifications" || !strings.Contains(body, `"on_failure":true`) {
		t.Fatalf("method=%q path=%q body=%q err=%v", method, path, body, err)
	}
	_, err = client.ListNotificationDeliveries(context.Background(), domain.NotificationDeliveryFilter{ChannelID: "channel one", TaskID: "task", State: domain.NotificationDeliveryFailed, Limit: 7})
	if err != nil || query != "channel=channel+one&limit=7&state=failed&task=task" {
		t.Fatalf("query=%q err=%v", query, err)
	}
}
