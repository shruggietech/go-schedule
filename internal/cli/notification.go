package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func newNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "notification", Short: "Manage run-outcome notifications"}
	cmd.AddCommand(newNotificationChannelCmd(), newNotificationTaskCmd(), newNotificationGroupCmd(), newNotificationDeliveriesCmd())
	return cmd
}

func newNotificationChannelCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "channel", Short: "Manage reusable webhook channels"}
	cmd.AddCommand(notificationChannelAdd(), notificationChannelList(), notificationChannelGet(), notificationChannelUpdate(), notificationChannelToggle(true), notificationChannelToggle(false), notificationChannelRotate(), notificationChannelTest(), notificationChannelRemove())
	return cmd
}

func notificationChannelAdd() *cobra.Command {
	var endpoint, authorization string
	var disabled bool
	cmd := &cobra.Command{Use: "add <name>", Short: "Create a webhook channel", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		if endpoint == "" {
			return fmtUsage("--endpoint is required")
		}
		enabled := !disabled
		ctx, cancel := reqCtx()
		defer cancel()
		channel, err := newClient().CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: args[0], Endpoint: endpoint, Authorization: authorization, Enabled: &enabled})
		if err != nil {
			return err
		}
		return printNotificationChannel(channel)
	}}
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "HTTPS webhook URL (HTTP allowed for loopback)")
	cmd.Flags().StringVar(&authorization, "authorization", "", "write-only Authorization header value")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create disabled")
	return cmd
}

func notificationChannelList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List webhook channels", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		channels, err := newClient().ListNotificationChannels(ctx)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(channels)
		}
		for _, channel := range channels {
			fmt.Fprintf(os.Stdout, "%s\t%s\t%s\tenabled=%t\tauthorization=%t\n", channel.ID, channel.Name, channel.EndpointSummary, channel.Enabled, channel.HasAuthorization)
		}
		return nil
	}}
}
func notificationChannelGet() *cobra.Command {
	return &cobra.Command{Use: "get <id>", Short: "Show a webhook channel", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		channel, err := newClient().GetNotificationChannel(ctx, args[0])
		if err != nil {
			return err
		}
		return printNotificationChannel(channel)
	}}
}

func notificationChannelUpdate() *cobra.Command {
	var name, endpoint string
	cmd := &cobra.Command{Use: "update <id>", Short: "Update a webhook channel", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		req := server.NotificationChannelUpdateRequest{}
		if cmd.Flags().Changed("name") {
			req.Name = &name
		}
		if cmd.Flags().Changed("endpoint") {
			req.Endpoint = &endpoint
		}
		if req.Name == nil && req.Endpoint == nil {
			return fmtUsage("at least one of --name or --endpoint is required")
		}
		ctx, cancel := reqCtx()
		defer cancel()
		channel, err := newClient().UpdateNotificationChannel(ctx, args[0], req)
		if err != nil {
			return err
		}
		return printNotificationChannel(channel)
	}}
	cmd.Flags().StringVar(&name, "name", "", "new channel name")
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "new webhook URL")
	return cmd
}

func notificationChannelToggle(enabled bool) *cobra.Command {
	verb := "disable"
	short := "Disable a webhook channel"
	if enabled {
		verb = "enable"
		short = "Enable a webhook channel"
	}
	return &cobra.Command{Use: verb + " <id>", Short: short, Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		channel, err := newClient().SetNotificationChannelEnabled(ctx, args[0], enabled)
		if err != nil {
			return err
		}
		return printNotificationChannel(channel)
	}}
}

func notificationChannelRotate() *cobra.Command {
	var authorization string
	cmd := &cobra.Command{Use: "rotate <id>", Short: "Replace or clear channel authorization", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		channel, err := newClient().RotateNotificationChannelAuthorization(ctx, args[0], authorization)
		if err != nil {
			return err
		}
		return printNotificationChannel(channel)
	}}
	cmd.Flags().StringVar(&authorization, "authorization", "", "new write-only Authorization value (empty clears)")
	return cmd
}

func notificationChannelTest() *cobra.Command {
	return &cobra.Command{Use: "test <id>", Short: "Queue a test webhook", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		delivery, err := newClient().TestNotificationChannel(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(delivery)
		}
		fmt.Fprintf(os.Stdout, "queued test delivery %s for channel %s\n", delivery.ID, delivery.ChannelName)
		return nil
	}}
}
func notificationChannelRemove() *cobra.Command {
	return &cobra.Command{Use: "rm <id>", Short: "Remove a webhook channel", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		if err := newClient().DeleteNotificationChannel(ctx, args[0]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "deleted notification channel %s\n", args[0])
		return nil
	}}
}

func printNotificationChannel(channel domain.NotificationChannel) error {
	if jsonOut {
		return printJSON(channel)
	}
	fmt.Fprintf(os.Stdout, "%s (%s) %s enabled=%t authorization=%t\n", channel.Name, channel.ID, channel.EndpointSummary, channel.Enabled, channel.HasAuthorization)
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
	var outcomes string
	cmd := &cobra.Command{Use: "set <id>", Short: "Replace the complete direct notification policy", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		onSuccess, onFailure, err := parseNotificationOutcomes(outcomes)
		if err != nil {
			return err
		}
		items := make([]server.NotificationAssignmentInput, len(channels))
		for i, id := range channels {
			items[i] = server.NotificationAssignmentInput{ChannelID: id, OnSuccess: onSuccess, OnFailure: onFailure}
		}
		ctx, cancel := reqCtx()
		defer cancel()
		var assignments []domain.NotificationAssignment
		if taskScope {
			assignments, err = newClient().ReplaceTaskNotificationAssignments(ctx, args[0], server.NotificationAssignmentsRequest{Assignments: items})
		} else {
			assignments, err = newClient().ReplaceGroupNotificationAssignments(ctx, args[0], server.NotificationAssignmentsRequest{Assignments: items})
		}
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(assignments)
		}
		fmt.Fprintf(os.Stdout, "configured %d notification assignment(s)\n", len(assignments))
		return nil
	}}
	cmd.Flags().StringSliceVar(&channels, "channel", nil, "channel ID (repeatable); omit all to resume inheritance")
	cmd.Flags().StringVar(&outcomes, "on", "failure", "comma-separated outcomes: success,failure")
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
		case "":
		default:
			return false, false, fmtUsage("--on accepts success and failure")
		}
	}
	if !success && !failure {
		return false, false, fmtUsage("--on must select success or failure")
	}
	return success, failure, nil
}

func notificationScopeShow(taskScope bool) *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show directly configured notification policy", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		var assignments []domain.NotificationAssignment
		var err error
		if taskScope {
			assignments, err = newClient().ListTaskNotificationAssignments(ctx, args[0])
		} else {
			assignments, err = newClient().ListGroupNotificationAssignments(ctx, args[0])
		}
		if err != nil {
			return err
		}
		return printJSON(assignments)
	}}
}
func notificationTaskEffective() *cobra.Command {
	return &cobra.Command{Use: "effective <id>", Short: "Explain the effective task notification policy", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		policy, err := newClient().EffectiveTaskNotificationPolicy(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(policy)
		}
		fmt.Fprintf(os.Stdout, "source=%s %s assignments=%d\n", policy.SourceScopeType, policy.SourceScopeID, len(policy.Assignments))
		return nil
	}}
}

func newNotificationDeliveriesCmd() *cobra.Command {
	var channelID, taskID, runID, state string
	var limit int
	cmd := &cobra.Command{Use: "deliveries", Short: "List redacted delivery evidence", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		deliveries, err := newClient().ListNotificationDeliveries(ctx, domain.NotificationDeliveryFilter{ChannelID: channelID, TaskID: taskID, RunID: runID, State: domain.NotificationDeliveryState(state), Limit: limit})
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(deliveries)
		}
		for _, delivery := range deliveries {
			fmt.Fprintf(os.Stdout, "%s\t%s\t%s\tattempts=%d\tstatus=%d\n", delivery.ID, delivery.EventKind, delivery.State, delivery.Attempts, delivery.LastStatus)
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
