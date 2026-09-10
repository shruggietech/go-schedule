package remote

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/config"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
)

type actorContextKey struct{}

func ActorID(r *http.Request) (string, error) {
	value, _ := r.Context().Value(actorContextKey{}).(string)
	if value == "" {
		return "", errors.New("remote actor unavailable")
	}
	return value, nil
}

type Handler struct {
	next        http.Handler
	enrollment  *enrollment.Service
	sources     *limiterSet
	actors      *limiterSet
	concurrent  chan struct{}
	streams     *streamSet
	publicActor string
}

func NewHandler(next http.Handler, enrollmentService *enrollment.Service, publicActorID string) *Handler {
	return &Handler{next: next, enrollment: enrollmentService, sources: newLimiterSet(10, 20, 4096), actors: newLimiterSet(20, 40, 4096), concurrent: make(chan struct{}, 64), streams: newStreamSet(16, 2), publicActor: publicActorID}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Header.Get("Origin") != "" {
		remoteError(w, http.StatusForbidden, "origin_rejected", "browser-origin requests are not supported")
		return
	}
	operation, ok := lookup(r.Method, r.URL.Path)
	if !ok {
		remoteError(w, http.StatusNotFound, "not_found", "remote endpoint not found")
		return
	}
	if operation.Retry != RetryReconnect {
		select {
		case h.concurrent <- struct{}{}:
			defer func() { <-h.concurrent }()
		default:
			remoteError(w, http.StatusServiceUnavailable, "busy", "remote request capacity is exhausted")
			return
		}
	}
	source, _, _ := net.SplitHostPort(r.RemoteAddr)
	if source == "" {
		source = r.RemoteAddr
	}
	if !h.sources.allow(source) {
		limited(w)
		return
	}
	if operation.BodyLimit > 0 {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			remoteError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "application/json is required")
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, operation.BodyLimit+1))
		if err != nil {
			remoteError(w, http.StatusBadRequest, "invalid_body", "request body could not be read")
			return
		}
		if int64(len(body)) > operation.BodyLimit {
			remoteError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds the operation limit")
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	if operation.ID == "enrollment.exchange" {
		h.exchange(w, r)
		return
	}
	if operation.Public {
		h.forward(w, r, operation, h.publicActor)
		return
	}
	token, ok := bearer(r.Header.Values("Authorization"))
	if !ok || r.URL.Query().Has("access_token") || r.Header.Get("Cookie") != "" {
		unauthorized(w)
		return
	}
	credential, actor, err := h.enrollment.Authenticate(token)
	if errors.Is(err, enrollment.ErrRevoked) {
		revoked(w)
		return
	}
	if err != nil || !h.actors.allow(credential.ID) {
		if err == nil {
			limited(w)
		} else {
			unauthorized(w)
		}
		return
	}
	if operation.Retry == RetryReconnect {
		if !h.streams.acquire(credential.ID) {
			limited(w)
			return
		}
		defer h.streams.release(credential.ID)
		ctx, cancel := context.WithCancel(r.Context())
		done := make(chan struct{})
		go func() {
			defer close(done)
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if _, _, err := h.enrollment.Authenticate(token); err != nil {
						cancel()
						return
					}
				}
			}
		}()
		h.forward(w, r.WithContext(ctx), operation, actor.ID)
		cancel()
		<-done
		return
	}
	h.forward(w, r, operation, actor.ID)
}

func (h *Handler) exchange(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PairingID   string            `json:"pairing_id"`
		Phrase      string            `json:"phrase"`
		DaemonID    string            `json:"daemon_id"`
		DisplayName string            `json:"display_name"`
		Kind        domain.ActorKind  `json:"kind"`
		Capability  domain.Capability `json:"capability"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&request) != nil || decoder.Decode(&struct{}{}) == nil {
		unauthorized(w)
		return
	}
	result, err := h.enrollment.Exchange(request.PairingID, request.Phrase, request.DaemonID, request.DisplayName, request.Kind, request.Capability)
	if err != nil {
		unauthorized(w)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) forward(w http.ResponseWriter, r *http.Request, operation Operation, actorID string) {
	clone := r.Clone(context.WithValue(r.Context(), actorContextKey{}, actorID))
	clone.URL.Path = localPath(operation, r.URL.Path)
	if operation.ID == "tasks.list" || operation.ID == "tasks.create" || operation.ID == "tasks.read" || operation.ID == "tasks.update" || operation.ID == "events.stream" {
		query := clone.URL.Query()
		query.Set("observation", "true")
		query.Del("details")
		clone.URL.RawQuery = query.Encode()
	}
	clone.Header.Set("X-Go-Schedule-Transport", "remote")
	h.next.ServeHTTP(w, clone)
}

func bearer(values []string) (string, bool) {
	if len(values) != 1 {
		return "", false
	}
	parts := strings.Split(values[0], " ")
	returnValue := ""
	valid := len(parts) == 2 && parts[0] == "Bearer" && parts[1] != ""
	if valid {
		returnValue = parts[1]
	}
	return returnValue, valid
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	remoteError(w, http.StatusUnauthorized, "unauthorized", "authentication failed")
}
func revoked(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	remoteError(w, http.StatusUnauthorized, "credential_revoked", "credential revoked")
}
func limited(w http.ResponseWriter) {
	w.Header().Set("Retry-After", "1")
	remoteError(w, http.StatusTooManyRequests, "rate_limited", "request rate limit exceeded")
}
func remoteError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, server.APIError{Error: server.ErrorBody{Code: code, Message: message}})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Serve(ctx context.Context, cfg config.RemoteConfig, handler http.Handler, log *slog.Logger) error {
	if !cfg.Enabled {
		<-ctx.Done()
		return nil
	}
	certificate, err := tls.LoadX509KeyPair(cfg.CertificateFile, cfg.PrivateKeyFile)
	if err != nil {
		return errors.New("remote: load TLS certificate: " + err.Error())
	}
	listener, err := net.Listen("tcp", cfg.BindAddress)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(listener, TLSConfig(certificate))
	httpServer := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32 << 10}
	done := make(chan error, 1)
	go func() {
		err := httpServer.Serve(tlsListener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		done <- err
	}()
	log.Info("remote HTTPS listener ready", "address", cfg.BindAddress)
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return errors.Join(httpServer.Shutdown(shutdown), <-done)
	}
}

func TLSConfig(certificate tls.Certificate) *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{certificate}, MinVersion: tls.VersionTLS13}
}
