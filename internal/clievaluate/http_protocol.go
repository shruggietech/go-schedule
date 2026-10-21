package clievaluate

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/domain"
)

type authorizedTransport struct {
	base  http.RoundTripper
	token string
}

func (t authorizedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	copy := request.Clone(request.Context())
	copy.Header = request.Header.Clone()
	copy.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(copy)
}

func (s *session) httpPost(ctx context.Context, definition step, body any) (bool, any, error) {
	status, err := s.api.MCPHTTPStatus(ctx)
	if err != nil {
		return false, status, err
	}
	switch definition.Spec.ID {
	case "mcp.http.status":
		return !status.Enabled, status, nil
	case "mcp.http.disable":
		return !status.Enabled, status, nil
	}
	value, ok := fieldAt(body, "credential")
	token, valid := value.(string)
	if !ok || !valid || token == "" {
		return false, status, errors.New("one-time localhost credential is absent")
	}
	s.redactor.Register(token)
	parsed, err := url.Parse(status.Endpoint)
	if err != nil || parsed.Hostname() != "127.0.0.1" || parsed.Port() != s.vars["mcp_port"] || parsed.Scheme != "http" {
		return false, status, errors.New("localhost MCP endpoint differs from the owned listener")
	}
	good := status.Enabled && status.Permission == domain.CapabilityManage && status.RequireConfirmation
	old := s.vars["mcp_token"]
	s.vars["mcp_token"] = token
	transport := http.DefaultTransport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: authorizedTransport{transport, token}, Timeout: 3 * time.Second}
	sdk := mcp.NewClient(&mcp.Implementation{Name: "QA localhost", Version: Version}, nil)
	connection, err := sdk.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: status.Endpoint, HTTPClient: client}, nil)
	var listed any
	var mutations any
	if err == nil {
		resources, e := connection.ListResources(ctx, nil)
		listed = resources
		err = errors.Join(err, e)
		good = good && e == nil && len(resources.Resources) == 5
		if err == nil {
			passed, details, e := s.manageProtocol(ctx, connection)
			good = good && passed
			mutations = details
			err = errors.Join(err, e)
		}
		err = errors.Join(err, connection.Close())
	}
	oldRejected := true
	if definition.Spec.ID == "mcp.http.rotate" {
		request, e := http.NewRequestWithContext(ctx, http.MethodPost, status.Endpoint, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"QA stale token","version":"1"}}}`))
		if e != nil {
			return false, status, e
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json, text/event-stream")
		request.Header.Set("Authorization", "Bearer "+old)
		response, e := transport.RoundTrip(request)
		oldRejected = e == nil && response.StatusCode == http.StatusUnauthorized
		if response != nil {
			e = errors.Join(e, response.Body.Close())
		}
		err = errors.Join(err, e)
	}
	return good && oldRejected, map[string]any{"status": status, "resources": listed, "mutations": mutations, "rotated_credential_rejected": oldRejected}, err
}
