package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/notification"
	"github.com/shruggietech/go-schedule/internal/store"
)

// NotificationChannelCreateRequest is the write-only channel creation contract.
type NotificationChannelCreateRequest struct {
	Name                  string `json:"name"`
	Endpoint              string `json:"endpoint"`
	Authorization         string `json:"authorization,omitempty"`
	Enabled               *bool  `json:"enabled,omitempty"`
	HealthIntervalSeconds int64  `json:"health_interval_seconds,omitempty"`
}

// NotificationChannelUpdateRequest changes non-authorization channel fields.
type NotificationChannelUpdateRequest struct {
	Name                  *string `json:"name,omitempty"`
	Endpoint              *string `json:"endpoint,omitempty"`
	Enabled               *bool   `json:"enabled,omitempty"`
	HealthIntervalSeconds *int64  `json:"health_interval_seconds,omitempty"`
}

// NotificationChannelRotateRequest replaces or removes authorization.
type NotificationChannelRotateRequest struct {
	Authorization string `json:"authorization"`
}

// NotificationAssignmentInput is one item in a complete scope replacement.
type NotificationAssignmentInput struct {
	ChannelID                string `json:"channel_id"`
	OnSuccess                bool   `json:"on_success"`
	OnFailure                bool   `json:"on_failure"`
	FailureThreshold         int    `json:"failure_threshold,omitempty"`
	OnFailureToStart         bool   `json:"on_failure_to_start,omitempty"`
	DurationThresholdSeconds int64  `json:"duration_threshold_seconds,omitempty"`
	OnRecovery               bool   `json:"on_recovery,omitempty"`
	ReminderIntervalSeconds  int64  `json:"reminder_interval_seconds,omitempty"`
	QuietPeriodSeconds       int64  `json:"quiet_period_seconds,omitempty"`
}

// NotificationAssignmentsRequest atomically replaces a scope policy.
type NotificationAssignmentsRequest struct {
	Assignments []NotificationAssignmentInput `json:"assignments"`
}

func (s *Server) handleCreateNotificationChannel(w http.ResponseWriter, r *http.Request) {
	var req NotificationChannelCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, CodeValidation, "name", "name is required")
		return
	}
	summary, err := notification.ValidateEndpoint(req.Endpoint)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "endpoint", err.Error())
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	if !validNotificationInterval(req.HealthIntervalSeconds, 60, 86400) {
		writeError(w, http.StatusBadRequest, CodeValidation, "health_interval_seconds", "health interval must be 0 or between 60 and 86400 seconds")
		return
	}
	channel := domain.NotificationChannel{Name: name, Kind: domain.NotificationChannelWebhook, Endpoint: strings.TrimSpace(req.Endpoint), EndpointSummary: summary, Authorization: req.Authorization, Enabled: enabled, HealthIntervalSeconds: req.HealthIntervalSeconds}
	if err := s.store.CreateNotificationChannel(&channel); err != nil {
		s.internalError(w, "create notification channel", err)
		return
	}
	writeJSON(w, http.StatusCreated, channel)
}

func (s *Server) handleListNotificationChannels(w http.ResponseWriter, _ *http.Request) {
	channels, err := s.store.ListNotificationChannels()
	if err != nil {
		s.internalError(w, "list notification channels", err)
		return
	}
	if channels == nil {
		channels = []domain.NotificationChannel{}
	}
	writeJSON(w, http.StatusOK, struct {
		Channels []domain.NotificationChannel `json:"notification_channels"`
	}{channels})
}

func (s *Server) handleGetNotificationChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := s.store.GetNotificationChannel(r.PathValue("id"))
	if s.notificationStoreError(w, err, "get notification channel") {
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (s *Server) handleUpdateNotificationChannel(w http.ResponseWriter, r *http.Request) {
	channel, err := s.store.GetNotificationChannel(r.PathValue("id"))
	if s.notificationStoreError(w, err, "get notification channel") {
		return
	}
	var req NotificationChannelUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON body")
		return
	}
	if req.Name != nil {
		channel.Name = strings.TrimSpace(*req.Name)
		if channel.Name == "" {
			writeError(w, http.StatusBadRequest, CodeValidation, "name", "name is required")
			return
		}
	}
	if req.Endpoint != nil {
		summary, err := notification.ValidateEndpoint(*req.Endpoint)
		if err != nil {
			writeError(w, http.StatusBadRequest, CodeValidation, "endpoint", err.Error())
			return
		}
		channel.Endpoint, channel.EndpointSummary = strings.TrimSpace(*req.Endpoint), summary
	}
	if req.Enabled != nil {
		channel.Enabled = *req.Enabled
	}
	if req.HealthIntervalSeconds != nil {
		if !validNotificationInterval(*req.HealthIntervalSeconds, 60, 86400) {
			writeError(w, http.StatusBadRequest, CodeValidation, "health_interval_seconds", "health interval must be 0 or between 60 and 86400 seconds")
			return
		}
		channel.HealthIntervalSeconds = *req.HealthIntervalSeconds
	}
	if err := s.store.UpdateNotificationChannel(channel); err != nil {
		s.internalError(w, "update notification channel", err)
		return
	}
	channel, err = s.store.GetNotificationChannel(channel.ID)
	if s.notificationStoreError(w, err, "get updated notification channel") {
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (s *Server) handleDeleteNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if s.notificationStoreError(w, s.store.DeleteNotificationChannel(r.PathValue("id")), "delete notification channel") {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) handleEnableNotificationChannel(w http.ResponseWriter, r *http.Request) {
	s.handleSetNotificationChannelEnabled(w, r, true)
}
func (s *Server) handleDisableNotificationChannel(w http.ResponseWriter, r *http.Request) {
	s.handleSetNotificationChannelEnabled(w, r, false)
}

func (s *Server) handleSetNotificationChannelEnabled(w http.ResponseWriter, r *http.Request, enabled bool) {
	id := r.PathValue("id")
	if s.notificationStoreError(w, s.store.SetNotificationChannelEnabled(id, enabled), "set notification channel state") {
		return
	}
	channel, err := s.store.GetNotificationChannel(id)
	if s.notificationStoreError(w, err, "get updated notification channel") {
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (s *Server) handleRotateNotificationChannel(w http.ResponseWriter, r *http.Request) {
	var req NotificationChannelRotateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON body")
		return
	}
	id := r.PathValue("id")
	if s.notificationStoreError(w, s.store.RotateNotificationChannelAuthorization(id, req.Authorization), "rotate notification authorization") {
		return
	}
	channel, err := s.store.GetNotificationChannel(id)
	if s.notificationStoreError(w, err, "get rotated notification channel") {
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (s *Server) handleTestNotificationChannel(w http.ResponseWriter, r *http.Request) {
	delivery, err := s.store.CreateTestNotificationDelivery(r.PathValue("id"))
	if errors.Is(err, store.ErrNotificationChannelDisabled) {
		writeError(w, http.StatusConflict, CodeConflict, "", "notification channel is disabled")
		return
	}
	if s.notificationStoreError(w, err, "create test notification") {
		return
	}
	if s.notify != nil {
		s.notify.Wake()
	}
	writeJSON(w, http.StatusAccepted, delivery)
}

func (s *Server) handleTaskNotifications(w http.ResponseWriter, r *http.Request) {
	s.handleScopeNotifications(w, r, domain.NotificationScopeTask)
}
func (s *Server) handleGroupNotifications(w http.ResponseWriter, r *http.Request) {
	s.handleScopeNotifications(w, r, domain.NotificationScopeGroup)
}

func (s *Server) handleScopeNotifications(w http.ResponseWriter, r *http.Request, scope domain.NotificationScopeType) {
	id := r.PathValue("id")
	if scope == domain.NotificationScopeTask {
		if _, err := s.store.GetTask(id); s.notificationStoreError(w, err, "get notification task") {
			return
		}
	} else {
		if _, err := s.store.GetGroup(id); s.notificationStoreError(w, err, "get notification group") {
			return
		}
	}
	if r.Method == http.MethodGet {
		assignments, err := s.store.ListNotificationAssignments(scope, id)
		if err != nil {
			s.internalError(w, "list notification assignments", err)
			return
		}
		if assignments == nil {
			assignments = []domain.NotificationAssignment{}
		}
		writeJSON(w, http.StatusOK, struct {
			Assignments []domain.NotificationAssignment `json:"assignments"`
		}{assignments})
		return
	}
	var req NotificationAssignmentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON body")
		return
	}
	assignments := make([]domain.NotificationAssignment, len(req.Assignments))
	seen := map[string]bool{}
	for i, item := range req.Assignments {
		if seen[item.ChannelID] {
			writeError(w, http.StatusConflict, CodeConflict, "assignments", "channel may be assigned only once per scope")
			return
		}
		seen[item.ChannelID] = true
		if item.FailureThreshold == 0 {
			item.FailureThreshold = 1
		}
		if item.FailureThreshold < 1 || item.FailureThreshold > 100 {
			writeError(w, http.StatusBadRequest, CodeValidation, "assignments", "failure threshold must be between 1 and 100")
			return
		}
		if !validNotificationInterval(item.DurationThresholdSeconds, 1, 2592000) || !validNotificationInterval(item.ReminderIntervalSeconds, 60, 2592000) || !validNotificationInterval(item.QuietPeriodSeconds, 60, 2592000) {
			writeError(w, http.StatusBadRequest, CodeValidation, "assignments", "duration must be 0 or between 1 and 2592000 seconds; reminder and quiet intervals must be 0 or between 60 and 2592000 seconds")
			return
		}
		assignments[i] = domain.NotificationAssignment{ChannelID: item.ChannelID, OnSuccess: item.OnSuccess, OnFailure: item.OnFailure, FailureThreshold: item.FailureThreshold, OnFailureToStart: item.OnFailureToStart, DurationThresholdSeconds: item.DurationThresholdSeconds, OnRecovery: item.OnRecovery, ReminderIntervalSeconds: item.ReminderIntervalSeconds, QuietPeriodSeconds: item.QuietPeriodSeconds}
	}
	if err := s.store.ReplaceNotificationAssignments(scope, id, assignments); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "", "scope or channel not found")
		} else {
			writeError(w, http.StatusBadRequest, CodeValidation, "assignments", err.Error())
		}
		return
	}
	assignments, err := s.store.ListNotificationAssignments(scope, id)
	if err != nil {
		s.internalError(w, "list replaced notification assignments", err)
		return
	}
	if assignments == nil {
		assignments = []domain.NotificationAssignment{}
	}
	writeJSON(w, http.StatusOK, struct {
		Assignments []domain.NotificationAssignment `json:"assignments"`
	}{assignments})
}

func validNotificationInterval(value, minimum, maximum int64) bool {
	return value == 0 || (value >= minimum && value <= maximum)
}

func (s *Server) handleEffectiveTaskNotifications(w http.ResponseWriter, r *http.Request) {
	policy, err := s.store.EffectiveNotificationPolicy(r.PathValue("id"))
	if s.notificationStoreError(w, err, "get effective notification policy") {
		return
	}
	writeJSON(w, http.StatusOK, policy)
}

func (s *Server) handleListNotificationDeliveries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 100
	if raw := q.Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 1000 {
			writeError(w, http.StatusBadRequest, CodeValidation, "limit", "limit must be between 1 and 1000")
			return
		}
		limit = parsed
	}
	state := domain.NotificationDeliveryState(q.Get("state"))
	if state != "" && state != domain.NotificationDeliveryPending && state != domain.NotificationDeliveryClaimed && state != domain.NotificationDeliverySucceeded && state != domain.NotificationDeliveryFailed {
		writeError(w, http.StatusBadRequest, CodeValidation, "state", "invalid delivery state")
		return
	}
	deliveries, err := s.store.ListNotificationDeliveries(domain.NotificationDeliveryFilter{ChannelID: q.Get("channel"), TaskID: q.Get("task"), RunID: q.Get("run"), State: state, Limit: limit})
	if err != nil {
		s.internalError(w, "list notification deliveries", err)
		return
	}
	if deliveries == nil {
		deliveries = []domain.NotificationDelivery{}
	}
	writeJSON(w, http.StatusOK, struct {
		Deliveries []domain.NotificationDelivery `json:"notification_deliveries"`
	}{deliveries})
}

func (s *Server) notificationStoreError(w http.ResponseWriter, err error, operation string) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, CodeNotFound, "", "notification resource not found")
	} else {
		s.internalError(w, operation, err)
	}
	return true
}

func (s *Server) internalError(w http.ResponseWriter, operation string, err error) {
	if s.log != nil {
		s.log.Error("api: "+operation, "err", err)
	}
	writeError(w, http.StatusInternalServerError, CodeInternal, "", operation+" failed")
}
