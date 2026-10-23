package clievaluate

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/clientsecret"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/userprefs"
)

func (s *session) localPost(definition step) (bool, any, error) {
	if strings.HasPrefix(definition.Spec.CommandPath, "preferences ") {
		prefs, err := userprefs.NewStore(userprefs.Options{Path: filepath.Join(s.w.Runtime, "client", "preferences.json"), LegacyPath: filepath.Join(s.w.Runtime, "client", "legacy.json")}).Load()
		if err != nil {
			return false, prefs, err
		}
		good := true
		switch definition.Spec.ID {
		case "preferences.show":
			good = prefs.Appearance == userprefs.AppearanceDark && prefs.Transition.Status == "migrated"
		case "preferences.set":
			good = prefs.Appearance == userprefs.AppearanceDark && prefs.Popups.Enabled && slices.Equal(prefs.Popups.Conditions, []string{"failure"}) && slices.Equal(prefs.Popups.Severities, []string{"error"}) && slices.Equal(prefs.Popups.DaemonIDs, []string{"fixture"})
		case "preferences.reset", "preferences.invalid-store-repair":
			good = prefs.Appearance == userprefs.AppearanceSystem && !prefs.Popups.Enabled
		case "preferences.reset-appearance":
			good = prefs.Appearance == userprefs.AppearanceSystem
		case "preferences.reset-popups":
			good = equalObservation(prefs.Popups, userprefs.DefaultPopups())
		}
		return good, prefs, nil
	}
	collection, err := s.profiles.Load()
	if err != nil {
		return false, collection, err
	}
	good := true
	switch definition.Spec.ID {
	case "profile.rename":
		profile, err := collection.Find(s.vars["profile"])
		good = err == nil && profile.Label == "QA renamed profile"
	case "profile.default.set":
		good = collection.ActiveDesktopProfileID == s.vars["profile"]
	case "profile.default.reset":
		good = collection.ActiveDesktopProfileID == ""
	case "profile.remove":
		_, findErr := collection.Find(s.vars["profile"])
		_, secretErr := clientsecret.NewWithBackend(s.recorder).Load(s.w.DaemonID, s.vars["credential"])
		good = findErr != nil && secretErr != nil
	}
	return good, collection, nil
}

func (s *session) assignmentPost(ctx context.Context, definition step) (bool, any, error) {
	group := definition.Spec.CommandPath == "notification group set"
	id := s.vars["base_task"]
	if group {
		id = s.vars["group"]
		if definition.Hook == "group-policy" {
			id = s.vars["base_group"]
		}
	}
	var actual []domain.NotificationAssignment
	var err error
	if group {
		actual, err = s.api.ListGroupNotificationAssignments(ctx, id)
	} else {
		actual, err = s.api.ListTaskNotificationAssignments(ctx, id)
	}
	if err != nil {
		return false, actual, err
	}
	expected := []server.NotificationAssignmentInput{}
	if definition.Hook != "task-inherit" && (definition.Spec.ID == "notification.task-set" || definition.Spec.ID == "notification.group-stdin") {
		var request server.NotificationAssignmentsRequest
		if err := json.Unmarshal([]byte(s.vars["policy_json"]), &request); err != nil {
			return false, actual, err
		}
		expected = request.Assignments
	} else if definition.Hook != "task-inherit" {
		channel := s.vars["channel"]
		success := true
		if definition.Hook == "task-policy" {
			channel = s.vars["second_channel"]
			success = false
		}
		expected = []server.NotificationAssignmentInput{{ChannelID: channel, OnSuccess: success, OnFailure: true, OnFailureToStart: true, OnRecovery: true, FailureThreshold: 2, DurationThresholdSeconds: 1, QuietPeriodSeconds: 60, ReminderIntervalSeconds: 3600}}
		if !success {
			expected[0].FailureThreshold = 1
			expected[0].QuietPeriodSeconds = 0
		}
	}
	good := len(actual) == len(expected)
	for _, want := range expected {
		found := false
		for _, got := range actual {
			if got.ChannelID != want.ChannelID {
				continue
			}
			found = got.OnSuccess == want.OnSuccess && got.OnFailure == want.OnFailure && got.OnFailureToStart == want.OnFailureToStart && got.OnRecovery == want.OnRecovery && got.DurationThresholdSeconds == want.DurationThresholdSeconds && got.QuietPeriodSeconds == want.QuietPeriodSeconds && got.ReminderIntervalSeconds == want.ReminderIntervalSeconds
			if want.OnFailure {
				found = found && got.FailureThreshold == want.FailureThreshold
			}
		}
		good = good && found
	}
	return good, map[string]any{"expected": expected, "actual": actual}, nil
}
