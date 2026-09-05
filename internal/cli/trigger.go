package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func newTriggerCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "trigger", Short: "Invoke tasks from local external processes"}
	cmd.AddCommand(triggerCreate(), triggerList(), triggerShow(), triggerUpdate(), triggerSetEnabled(true), triggerSetEnabled(false), triggerRotate(), triggerRemove(), triggerFire())
	return cmd
}

func triggerCreate() *cobra.Command {
	var name, target string
	var disabled bool
	cmd := &cobra.Command{Use: "create", Short: "Create an external trigger", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		if name == "" || target == "" {
			return fmtUsage("--name and --task are required")
		}
		enabled := !disabled
		ctx, cancel := reqCtx()
		defer cancel()
		result, err := newClient().CreateTrigger(ctx, server.TriggerCreateRequest{Name: name, TargetTaskID: target, Enabled: &enabled})
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		fmt.Fprintf(os.Stdout, "created trigger %s\nkey: %s\ncommand: %s\n", result.Trigger.ID, result.Key, result.Command)
		return nil
	}}
	cmd.Flags().StringVar(&name, "name", "", "trigger name")
	cmd.Flags().StringVar(&target, "task", "", "target task ID")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create the trigger disabled")
	return cmd
}

func triggerList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List external triggers", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		items, err := newClient().ListTriggers(ctx)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(items)
		}
		tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tTARGET\tENABLED\tREADINESS")
		for _, item := range items {
			fmt.Fprintf(tw, "%s\t%s\t%s (%s)\t%t\t%s\n", item.ID, item.Name, item.TargetTaskName, item.TargetTaskID, item.Enabled, item.Readiness)
		}
		return tw.Flush()
	}}
}

func triggerShow() *cobra.Command {
	var reveal bool
	cmd := &cobra.Command{Use: "show <trigger-id>", Short: "Show an external trigger", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		if reveal {
			result, err := newClient().RevealTrigger(ctx, args[0])
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(result)
			}
			fmt.Fprintf(os.Stdout, "%s\nname: %s\ntarget: %s (%s)\nenabled: %t\nreadiness: %s\nkey: %s\ncommand: %s\n", result.Trigger.ID, result.Trigger.Name, result.Trigger.TargetTaskName, result.Trigger.TargetTaskID, result.Trigger.Enabled, result.Trigger.Readiness, result.Key, result.Command)
			return nil
		}
		item, err := newClient().GetTrigger(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(item)
		}
		fmt.Fprintf(os.Stdout, "%s\nname: %s\ntarget: %s (%s)\nenabled: %t\nreadiness: %s\n", item.ID, item.Name, item.TargetTaskName, item.TargetTaskID, item.Enabled, item.Readiness)
		return nil
	}}
	cmd.Flags().BoolVar(&reveal, "reveal-key", false, "include the sensitive trigger key")
	return cmd
}

func triggerUpdate() *cobra.Command {
	var name, target string
	cmd := &cobra.Command{Use: "update <trigger-id>", Short: "Update an external trigger", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		req := server.TriggerUpdateRequest{}
		if cmd.Flags().Changed("name") {
			if name == "" {
				return fmtUsage("--name cannot be empty")
			}
			req.Name = &name
		}
		if cmd.Flags().Changed("task") {
			if target == "" {
				return fmtUsage("--task cannot be empty")
			}
			req.TargetTaskID = &target
		}
		if req.Name == nil && req.TargetTaskID == nil {
			return fmtUsage("provide --name or --task")
		}
		ctx, cancel := reqCtx()
		defer cancel()
		item, err := newClient().UpdateTrigger(ctx, args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(item)
		}
		fmt.Fprintf(os.Stdout, "updated trigger %s\n", item.ID)
		return nil
	}}
	cmd.Flags().StringVar(&name, "name", "", "trigger name")
	cmd.Flags().StringVar(&target, "task", "", "target task ID")
	return cmd
}

func triggerSetEnabled(enabled bool) *cobra.Command {
	action, past := "disable", "disabled"
	if enabled {
		action, past = "enable", "enabled"
	}
	return &cobra.Command{Use: action + " <trigger-id>", Short: action + " an external trigger", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		item, err := newClient().SetTriggerEnabled(ctx, args[0], enabled)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(item)
		}
		fmt.Fprintf(os.Stdout, "%s trigger %s\n", past, item.ID)
		return nil
	}}
}

func triggerRotate() *cobra.Command {
	return &cobra.Command{Use: "rotate <trigger-id>", Short: "Replace an external trigger key", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		result, err := newClient().RotateTrigger(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		fmt.Fprintf(os.Stdout, "rotated trigger %s\nkey: %s\ncommand: %s\n", result.Trigger.ID, result.Key, result.Command)
		return nil
	}}
}

func triggerRemove() *cobra.Command {
	return &cobra.Command{Use: "rm <trigger-id>", Short: "Delete an external trigger", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		if err := newClient().DeleteTrigger(ctx, args[0]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "deleted trigger %s\n", args[0])
		return nil
	}}
}

func triggerFire() *cobra.Command {
	return &cobra.Command{Use: "fire <key>", Short: "Invoke an external trigger", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		if err := newClient().FireTrigger(ctx, args[0]); err != nil {
			return err
		}
		if jsonOut {
			return printJSON(map[string]string{"status": "accepted"})
		}
		fmt.Fprintln(os.Stdout, "trigger accepted")
		return nil
	}}
}
