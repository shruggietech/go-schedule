// Package mcphttp provides the optional authenticated loopback-only MCP transport.
package mcphttp

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxRequestBodyBytes = 1 << 20

func normalizeOrigins(values []string) ([]string, error) {
	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		u, err := url.Parse(value)
		if err != nil || u.Opaque != "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
			return nil, validationError("allowed_origins", "origins must be absolute loopback HTTP or HTTPS origins with no path, query, fragment, or credentials")
		}
		scheme := strings.ToLower(u.Scheme)
		if scheme != "http" && scheme != "https" || u.Hostname() != "127.0.0.1" || u.Port() == "" {
			return nil, validationError("allowed_origins", "origins must use http or https, numeric 127.0.0.1, and an explicit port")
		}
		port, err := strconv.Atoi(u.Port())
		if err != nil || port < 1 || port > 65535 {
			return nil, validationError("allowed_origins", "origin ports must be between 1 and 65535")
		}
		unique[scheme+"://127.0.0.1:"+strconv.Itoa(port)] = struct{}{}
	}
	origins := make([]string, 0, len(unique))
	for origin := range unique {
		origins = append(origins, origin)
	}
	sort.Strings(origins)
	return origins, nil
}

func streamableHandler(manager *Manager) http.Handler {
	sdk := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return manager.newObserveServer()
	}, &mcp.StreamableHTTPOptions{
		Stateless:                    true,
		MaxRequestBodyBytes:          maxRequestBodyBytes,
		PropagateRequestCancellation: true,
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		host, origins, digest, enabled := manager.requestPolicy()
		if !enabled || r.Host != host {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		requestOrigins := r.Header.Values("Origin")
		if len(requestOrigins) > 1 || len(requestOrigins) == 1 && !contains(origins, requestOrigins[0]) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		authorizations := r.Header.Values("Authorization")
		if len(authorizations) != 1 || !validBearer(authorizations[0], digest) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/mcp" {
			http.NotFound(w, r)
			return
		}
		sdk.ServeHTTP(w, r)
	})
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func validBearer(value string, expected [sha256.Size]byte) bool {
	if !strings.HasPrefix(value, "Bearer ") || strings.ContainsAny(value[len("Bearer "):], " \t\r\n") || len(value) == len("Bearer ") {
		return false
	}
	actual := sha256.Sum256([]byte(value[len("Bearer "):]))
	return subtle.ConstantTimeCompare(actual[:], expected[:]) == 1
}
