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
		if r.URL.Path != "/mcp" {
			http.NotFound(w, r)
			return
		}
		requestOrigins := r.Header.Values("Origin")
		if r.Method == http.MethodOptions {
			origin, ok := manager.authorizePreflight(r.Host, requestOrigins, r.Header.Get("Access-Control-Request-Method"), r.Header.Get("Access-Control-Request-Headers"))
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			setCORSHeaders(w, origin)
			w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, MCP-Protocol-Version")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		allowed, origin, unauthorized := manager.authorizeRequest(r.Host, requestOrigins, r.Header.Values("Authorization"))
		if !allowed && !unauthorized {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if unauthorized {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if origin != "" {
			setCORSHeaders(w, origin)
		}
		sdk.ServeHTTP(w, r)
	})
}

func (m *Manager) authorizeRequest(host string, origins, authorizations []string) (allowed bool, origin string, unauthorized bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.status.Enabled || host != endpointHost(m.status.Endpoint) || len(origins) > 1 || len(origins) == 1 && !contains(m.status.AllowedOrigins, origins[0]) {
		return false, "", false
	}
	if len(authorizations) != 1 || !validBearer(authorizations[0], m.digest) {
		return false, "", true
	}
	if len(origins) == 1 {
		origin = origins[0]
	}
	return true, origin, false
}

func (m *Manager) authorizePreflight(host string, origins []string, method, requestedHeaders string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.status.Enabled || host != endpointHost(m.status.Endpoint) || len(origins) != 1 || !contains(m.status.AllowedOrigins, origins[0]) || method != http.MethodPost {
		return "", false
	}
	allowedHeaders := map[string]bool{"authorization": true, "content-type": true, "accept": true, "mcp-protocol-version": true}
	for _, header := range strings.Split(requestedHeaders, ",") {
		header = strings.ToLower(strings.TrimSpace(header))
		if header != "" && !allowedHeaders[header] {
			return "", false
		}
	}
	return origins[0], true
}

func endpointHost(endpoint string) string {
	if !strings.HasPrefix(endpoint, "http://") || !strings.HasSuffix(endpoint, "/mcp") {
		return ""
	}
	return endpoint[len("http://") : len(endpoint)-len("/mcp")]
}

func setCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
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
