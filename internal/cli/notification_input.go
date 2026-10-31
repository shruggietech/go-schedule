package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

const notificationInputLimit = 1024 * 1024

func readNotificationInput(cmd *cobra.Command, path string) (server.NotificationAssignmentsRequest, error) {
	var request server.NotificationAssignmentsRequest
	if path == "" {
		return request, fmtUsage("--file requires a JSON path or - for standard input")
	}
	reader := cmd.InOrStdin()
	var file *os.File
	if path != "-" {
		var err error
		file, err = os.Open(path)
		if err != nil {
			return request, fmt.Errorf("open notification policy: %w", err)
		}
		reader = file
	}
	data, err := io.ReadAll(io.LimitReader(reader, notificationInputLimit+1))
	if file != nil {
		err = errors.Join(err, file.Close())
	}
	if err != nil {
		return request, fmt.Errorf("read notification policy: %w", err)
	}
	if len(data) > notificationInputLimit {
		return request, fmtUsage("notification policy exceeds 1 MiB; reduce the assignments request")
	}
	if !utf8.Valid(data) {
		return request, fmtUsage("notification policy must be UTF-8 text; save the JSON as UTF-8 and try again")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return request, fmtUsage(fmt.Sprintf("invalid notification policy JSON: %v", err))
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return request, fmtUsage("notification policy must contain exactly one JSON object with an assignments array")
	}
	if request.Assignments == nil {
		return request, fmtUsage("notification policy requires an assignments array; use {\"assignments\":[]} to resume inheritance")
	}
	return request, nil
}

func notificationInput(assignments []domain.NotificationAssignment) server.NotificationAssignmentsRequest {
	request := server.NotificationAssignmentsRequest{Assignments: make([]server.NotificationAssignmentInput, 0, len(assignments))}
	for _, item := range assignments {
		request.Assignments = append(request.Assignments, server.NotificationAssignmentInput{
			ChannelID: item.ChannelID, OnSuccess: item.OnSuccess, OnFailure: item.OnFailure,
			FailureThreshold: item.FailureThreshold, OnFailureToStart: item.OnFailureToStart,
			DurationThresholdSeconds: item.DurationThresholdSeconds, OnRecovery: item.OnRecovery,
			ReminderIntervalSeconds: item.ReminderIntervalSeconds, QuietPeriodSeconds: item.QuietPeriodSeconds,
		})
	}
	return request
}
