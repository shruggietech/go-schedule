package cli

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func newNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "notification", Short: "Manage run-outcome notifications"}
	cmd.AddCommand(newNotificationChannelCmd(), newNotificationTaskCmd(), newNotificationGroupCmd(), newNotificationDeliveriesCmd())
	return cmd
}

func newNotificationChannelCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "channel", Short: "Manage reusable notification channels"}
	cmd.AddCommand(notificationChannelAdd(), notificationChannelList(), notificationChannelGet(), notificationChannelUpdate(), notificationChannelToggle(true), notificationChannelToggle(false), notificationChannelRotate(), notificationChannelTest(), notificationChannelRemove())
	return cmd
}

func notificationChannelAdd() *cobra.Command {
	var endpoint, authorization, kind string
	var smtp smtpFlags
	var healthInterval time.Duration
	var disabled bool
	cmd := &cobra.Command{Use: "add <name>", Short: "Create a notification channel", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if kind == "webhook" && endpoint == "" {
			return fmtUsage("--endpoint is required for webhook channels")
		}
		if kind != "webhook" && kind != "smtp" {
			return fmtUsage("--kind accepts webhook or smtp")
		}
		if kind == "webhook" && (smtpSettingsChanged(cmd) || cmd.Flags().Changed("smtp-auth-file")) {
			return fmtUsage("SMTP flags require --kind smtp")
		}
		if kind == "smtp" && (cmd.Flags().Changed("endpoint") || cmd.Flags().Changed("authorization")) {
			return fmtUsage("SMTP channels use --smtp-* settings and --smtp-auth-file instead of webhook fields")
		}
		var settings *domain.SMTPConfig
		if kind == "smtp" {
			var err error
			settings, err = smtp.settings(cmd, nil)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("smtp-auth-file") {
				authorization, err = readSMTPAuth(cmd, smtp.authFile)
				if err != nil {
					return err
				}
			}
		}
		enabled := !disabled
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		channel, err := clientFor(cmd).CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: args[0], Kind: domain.NotificationChannelKind(kind), SMTP: settings, Endpoint: endpoint, Authorization: authorization, Enabled: &enabled, HealthIntervalSeconds: int64(healthInterval / time.Second)})
		if err != nil {
			return err
		}
		return printNotificationChannel(cmd, channel)
	}}
	cmd.Flags().StringVar(&kind, "kind", "webhook", "notification transport: webhook or smtp")
	smtp.bind(cmd, true)
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "HTTPS webhook URL (HTTP allowed for loopback)")
	cmd.Flags().StringVar(&authorization, "authorization", "", "write-only Authorization header value")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create disabled")
	cmd.Flags().DurationVar(&healthInterval, "health-interval", 0, "healthy-presence heartbeat interval (0 disables)")
	return cmd
}

func notificationChannelList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List notification channels", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		channels, err := clientFor(cmd).ListNotificationChannels(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), channels)
		}
		return writeHumanTable(cmd.OutOrStdout(), "CHANNEL ID\tNAME\tENDPOINT\tENABLED\tHAS AUTHORIZATION", func(out io.Writer) error {
			for _, channel := range channels {
				if _, err := fmt.Fprintf(out, "%s\t%s\t%s\t%t\t%t\n", channel.ID, channel.Name, channel.EndpointSummary, channel.Enabled, channel.HasAuthorization); err != nil {
					return err
				}
			}
			return nil
		})
	}}
}
func notificationChannelGet() *cobra.Command {
	return &cobra.Command{Use: "get <id>", Short: "Show a notification channel", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		channel, err := clientFor(cmd).GetNotificationChannel(ctx, args[0])
		if err != nil {
			return err
		}
		return printNotificationChannel(cmd, channel)
	}}
}

func notificationChannelUpdate() *cobra.Command {
	var name, endpoint string
	var smtp smtpFlags
	var healthInterval time.Duration
	cmd := &cobra.Command{Use: "update <id>", Short: "Update a notification channel", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		req := server.NotificationChannelUpdateRequest{}
		if cmd.Flags().Changed("name") {
			req.Name = &name
		}
		if cmd.Flags().Changed("endpoint") {
			req.Endpoint = &endpoint
		}
		if cmd.Flags().Changed("health-interval") {
			seconds := int64(healthInterval / time.Second)
			req.HealthIntervalSeconds = &seconds
		}
		if smtpSettingsChanged(cmd) {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			current, err := clientFor(cmd).GetNotificationChannel(ctx, args[0])
			if err != nil {
				return err
			}
			if current.Kind != domain.NotificationChannelSMTP {
				return fmtUsage("SMTP flags require an SMTP channel")
			}
			settings, err := smtp.settings(cmd, current.SMTP)
			if err != nil {
				return err
			}
			req.SMTP = settings
		}
		if req.Name == nil && req.Endpoint == nil && req.HealthIntervalSeconds == nil && req.SMTP == nil {
			return fmtUsage("provide --name, --endpoint, --health-interval, or an SMTP setting")
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		channel, err := clientFor(cmd).UpdateNotificationChannel(ctx, args[0], req)
		if err != nil {
			return err
		}
		return printNotificationChannel(cmd, channel)
	}}
	smtp.bind(cmd, false)
	cmd.Flags().StringVar(&name, "name", "", "new channel name")
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "new webhook URL")
	cmd.Flags().DurationVar(&healthInterval, "health-interval", 0, "healthy-presence heartbeat interval (0 disables)")
	return cmd
}

func notificationChannelToggle(enabled bool) *cobra.Command {
	verb := "disable"
	short := "Disable a notification channel"
	if enabled {
		verb = "enable"
		short = "Enable a notification channel"
	}
	return &cobra.Command{Use: verb + " <id>", Short: short, Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		channel, err := clientFor(cmd).SetNotificationChannelEnabled(ctx, args[0], enabled)
		if err != nil {
			return err
		}
		return printNotificationChannel(cmd, channel)
	}}
}

func notificationChannelRotate() *cobra.Command {
	var authorization, authFile string
	cmd := &cobra.Command{Use: "rotate <id>", Short: "Replace or clear channel authorization", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		if cmd.Flags().Changed("smtp-auth-file") {
			if cmd.Flags().Changed("authorization") {
				return fmtUsage("--smtp-auth-file cannot be combined with --authorization")
			}
			var err error
			authorization, err = readSMTPAuth(cmd, authFile)
			if err != nil {
				return err
			}
		}
		channel, err := clientFor(cmd).RotateNotificationChannelAuthorization(ctx, args[0], authorization)
		if err != nil {
			return err
		}
		return printNotificationChannel(cmd, channel)
	}}
	cmd.Flags().StringVar(&authFile, "smtp-auth-file", "", "write-only SMTP username/password JSON file; - reads standard input; {} clears")
	cmd.Flags().StringVar(&authorization, "authorization", "", "new write-only authorization value (empty clears)")
	return cmd
}

func notificationChannelTest() *cobra.Command {
	return &cobra.Command{Use: "test <id>", Short: "Queue a test notification", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		delivery, err := clientFor(cmd).TestNotificationChannel(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), delivery)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "queued test delivery %s for channel %s\n", delivery.ID, delivery.ChannelName)
		return nil
	}}
}
func notificationChannelRemove() *cobra.Command {
	return &cobra.Command{Use: "rm <id>", Short: "Remove a notification channel", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		if err := clientFor(cmd).DeleteNotificationChannel(ctx, args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted notification channel %s\n", args[0])
		return nil
	}}
}

func printNotificationChannel(cmd *cobra.Command, channel domain.NotificationChannel) error {
	if jsonOutput(cmd) {
		return printJSONTo(cmd.OutOrStdout(), channel)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s (%s) %s enabled=%t authorization=%t health_interval=%ds\n", channel.Name, channel.ID, channel.EndpointSummary, channel.Enabled, channel.HasAuthorization, channel.HealthIntervalSeconds)
	if channel.SMTP != nil {
		c := channel.SMTP
		fmt.Fprintf(cmd.OutOrStdout(), "SMTP host=%s port=%d encryption=%s from=%s recipients=%s timeout=%ds html=%t\n", c.Host, c.Port, c.TLSMode, c.From, strings.Join(c.Recipients, ","), c.TimeoutSeconds, c.HTML)
	}
	return nil
}

func newNotificationTaskCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "task", Short: "Manage task notification policy"}
	cmd.AddCommand(notificationScopeSet(true), notificationScopeShow(true), notificationTaskEffective())
	return cmd
}
func newNotificationGroupCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "group", Short: "Manage group notification policy"}
	cmd.AddCommand(notificationScopeSet(false), notificationScopeShow(false))
	return cmd
}

func notificationScopeSet(taskScope bool) *cobra.Command {
	var channels []string
	var filename string
	var outcomes string
	var failureThreshold int
	var failureToStart, recovery bool
	var durationThreshold, reminder, quietPeriod time.Duration
	cmd := &cobra.Command{Use: "set <id>", Short: "Replace the complete direct notification policy", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		var items []server.NotificationAssignmentInput
		if cmd.Flags().Changed("file") {
			for _, flag := range []string{"channel", "on", "failure-threshold", "failure-to-start", "duration-threshold", "recovery", "reminder", "quiet-period"} {
				if cmd.Flags().Changed(flag) {
					return fmtUsage("--file cannot be combined with --" + flag + "; put each destination's conditions in the JSON request")
				}
			}
			request, err := readNotificationInput(cmd, filename)
			if err != nil {
				return err
			}
			items = request.Assignments
		} else {
			onSuccess, onFailure, err := parseNotificationOutcomes(outcomes)
			if err != nil {
				return err
			}
			if !onSuccess && !onFailure && !failureToStart && durationThreshold == 0 && len(channels) > 0 {
				return fmtUsage("select at least one outcome, --failure-to-start, or --duration-threshold")
			}
			items = make([]server.NotificationAssignmentInput, len(channels))
			for i, id := range channels {
				items[i] = server.NotificationAssignmentInput{ChannelID: id, OnSuccess: onSuccess, OnFailure: onFailure, FailureThreshold: failureThreshold, OnFailureToStart: failureToStart, DurationThresholdSeconds: int64(durationThreshold / time.Second), OnRecovery: recovery, ReminderIntervalSeconds: int64(reminder / time.Second), QuietPeriodSeconds: int64(quietPeriod / time.Second)}
			}
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		var assignments []domain.NotificationAssignment
		var err error
		if taskScope {
			assignments, err = clientFor(cmd).ReplaceTaskNotificationAssignments(ctx, args[0], server.NotificationAssignmentsRequest{Assignments: items})
		} else {
			assignments, err = clientFor(cmd).ReplaceGroupNotificationAssignments(ctx, args[0], server.NotificationAssignmentsRequest{Assignments: items})
		}
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), assignments)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "configured %d notification assignment(s)\n", len(assignments))
		return nil
	}}
	cmd.Flags().StringSliceVar(&channels, "channel", nil, "channel ID (repeatable); omit all to resume inheritance")
	cmd.Flags().StringVar(&filename, "file", "", "UTF-8 JSON assignments request (at most 1 MiB); - reads standard input")
	cmd.Flags().StringVar(&outcomes, "on", "failure", "comma-separated outcomes: success,failure")
	cmd.Flags().IntVar(&failureThreshold, "failure-threshold", 1, "consecutive failures required before notification")
	cmd.Flags().BoolVar(&failureToStart, "failure-to-start", false, "notify when the task process cannot start")
	cmd.Flags().DurationVar(&durationThreshold, "duration-threshold", 0, "notify when run duration reaches this threshold")
	cmd.Flags().BoolVar(&recovery, "recovery", false, "notify when an active problem clears")
	cmd.Flags().DurationVar(&reminder, "reminder", 0, "repeat an active problem after this interval")
	cmd.Flags().DurationVar(&quietPeriod, "quiet-period", 0, "suppress routine notifications within this interval")
	return cmd
}

func parseNotificationOutcomes(value string) (bool, bool, error) {
	var success, failure bool
	for _, item := range strings.Split(value, ",") {
		switch strings.TrimSpace(item) {
		case "success":
			success = true
		case "failure":
			failure = true
		case "none":
		case "":
		default:
			return false, false, fmtUsage("--on accepts success, failure, and none")
		}
	}
	return success, failure, nil
}

func notificationScopeShow(taskScope bool) *cobra.Command {
	var input bool
	cmd := &cobra.Command{Use: "show <id>", Short: "Show directly configured notification policy", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		var assignments []domain.NotificationAssignment
		var err error
		if taskScope {
			assignments, err = clientFor(cmd).ListTaskNotificationAssignments(ctx, args[0])
		} else {
			assignments, err = clientFor(cmd).ListGroupNotificationAssignments(ctx, args[0])
		}
		if err != nil {
			return err
		}
		if input {
			return printJSONTo(cmd.OutOrStdout(), notificationInput(assignments))
		}
		return printJSONTo(cmd.OutOrStdout(), assignments)
	}}
	cmd.Flags().BoolVar(&input, "input", false, "export the canonical assignments request for editing and reapplying with set --file")
	return cmd
}
func notificationTaskEffective() *cobra.Command {
	return &cobra.Command{Use: "effective <id>", Short: "Explain the effective task notification policy", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		policy, err := clientFor(cmd).EffectiveTaskNotificationPolicy(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), policy)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "source=%s %s assignments=%d\n", policy.SourceScopeType, policy.SourceScopeID, len(policy.Assignments))
		for _, item := range policy.Assignments {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "destination %s: success=%t, failure=%t, failure threshold=%d, could not start=%t, duration threshold=%ds, recovery=%t, reminder=%ds, quiet period=%ds\n", item.ChannelID, item.OnSuccess, item.OnFailure, item.FailureThreshold, item.OnFailureToStart, item.DurationThresholdSeconds, item.OnRecovery, item.ReminderIntervalSeconds, item.QuietPeriodSeconds); err != nil {
				return err
			}
		}
		return nil
	}}
}

func newNotificationDeliveriesCmd() *cobra.Command {
	var channelID, taskID, runID, state string
	var limit int
	cmd := &cobra.Command{Use: "deliveries", Short: "List redacted delivery evidence", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		deliveries, err := clientFor(cmd).ListNotificationDeliveries(ctx, domain.NotificationDeliveryFilter{ChannelID: channelID, TaskID: taskID, RunID: runID, State: domain.NotificationDeliveryState(state), Limit: limit})
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), deliveries)
		}
		for _, delivery := range deliveries {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\tattempts=%d\tstatus=%d\n", delivery.ID, delivery.EventKind, delivery.State, delivery.Attempts, delivery.LastStatus)
			if delivery.Diagnostic != "" && delivery.Diagnostic != delivery.LastError {
				fmt.Fprintf(cmd.OutOrStdout(), "  Delivery detail: %s\n", delivery.Diagnostic)
			}
			if delivery.LastError != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Failure detail: %s\n", delivery.LastError)
			}
			if delivery.Kind == domain.NotificationChannelSMTP {
				for _, recipient := range delivery.Recipients {
					fmt.Fprintf(cmd.OutOrStdout(), "  recipient %s: %s, SMTP status=%d, %s\n", recipient.Address, recipient.State, recipient.Code, recipient.Diagnostic)
				}
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&channelID, "channel", "", "filter by channel ID")
	cmd.Flags().StringVar(&taskID, "task", "", "filter by task ID")
	cmd.Flags().StringVar(&runID, "run", "", "filter by run ID")
	cmd.Flags().StringVar(&state, "state", "", "filter by delivery state")
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum records (1-1000)")
	return cmd
}
