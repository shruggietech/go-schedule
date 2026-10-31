package cli

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

func newTriggerSetCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "set", Short: "Administer Trigger Sets"}
	cmd.AddCommand(triggerSetCreate(), triggerSetList(), triggerSetShow(), triggerSetRetarget(), triggerSetBulkEnabled(true), triggerSetBulkEnabled(false), triggerSetReveal(), triggerSetRotate(), triggerSetRemove())
	return cmd
}

func triggerSetCreate() *cobra.Command {
	var name, target string
	var count int
	var disabled bool
	cmd := &cobra.Command{Use: "create", Short: "Create a Trigger Set", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if name == "" || target == "" || count < 1 || count > 99 {
			return fmtUsage("--name, --task, and --count from 1 through 99 are required")
		}
		enabled := !disabled
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		result, err := clientFor(cmd).CreateTriggerSet(ctx, server.TriggerSetCreateRequest{Name: name, TargetTaskID: target, Count: count, Enabled: &enabled})
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), result)
		}
		return printTriggerSetSecret(cmd.OutOrStdout(), result)
	}}
	cmd.Flags().StringVar(&name, "name", "", "Trigger Set name")
	cmd.Flags().StringVar(&target, "task", "", "target task ID")
	cmd.Flags().IntVar(&count, "count", 0, "number of trigger members from 1 through 99")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create every member disabled")
	return cmd
}

func triggerSetList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List Trigger Sets", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		items, err := clientFor(cmd).ListTriggerSets(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), items)
		}
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tTARGET\tMEMBERS\tENABLED")
		for _, item := range items {
			fmt.Fprintf(tw, "%s\t%s\t%s (%s)\t%d\t%d\n", item.ID, item.Name, item.TargetTaskName, item.TargetTaskID, item.MemberCount, item.EnabledCount)
		}
		return tw.Flush()
	}}
}

func triggerSetShow() *cobra.Command {
	return &cobra.Command{Use: "show <set-id>", Short: "Show a Trigger Set", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := clientFor(cmd).GetTriggerSet(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), item)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\nname: %s\ntarget: %s (%s)\nmembers: %d\nenabled: %d\n", item.ID, item.Name, item.TargetTaskName, item.TargetTaskID, item.MemberCount, item.EnabledCount)
		return nil
	}}
}

func triggerSetReveal() *cobra.Command {
	return triggerSetSecretAction("reveal", "Reveal all Trigger Set commands", func(cmd *cobra.Command, id string) (server.TriggerSetSecretResponse, error) {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		return clientFor(cmd).RevealTriggerSet(ctx, id)
	})
}
func triggerSetRotate() *cobra.Command {
	return triggerSetSecretAction("rotate", "Rotate every Trigger Set key", func(cmd *cobra.Command, id string) (server.TriggerSetSecretResponse, error) {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		return clientFor(cmd).RotateTriggerSet(ctx, id)
	})
}

func triggerSetSecretAction(use, short string, action func(*cobra.Command, string) (server.TriggerSetSecretResponse, error)) *cobra.Command {
	return &cobra.Command{Use: use + " <set-id>", Short: short, Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		result, err := action(cmd, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), result)
		}
		return printTriggerSetSecret(cmd.OutOrStdout(), result)
	}}
}

func printTriggerSetSecret(out io.Writer, result server.TriggerSetSecretResponse) error {
	for _, member := range result.Members {
		if _, err := fmt.Fprintln(out, member.Command); err != nil {
			return err
		}
	}
	return nil
}

func triggerSetRetarget() *cobra.Command {
	var target string
	cmd := &cobra.Command{Use: "retarget <set-id>", Short: "Retarget every Trigger Set member", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if target == "" {
			return fmtUsage("--task is required")
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := clientFor(cmd).RetargetTriggerSet(ctx, args[0], target)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), item)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "retargeted Trigger Set %s (%d members)\n", item.ID, item.MemberCount)
		return nil
	}}
	cmd.Flags().StringVar(&target, "task", "", "new target task ID")
	return cmd
}

func triggerSetBulkEnabled(enabled bool) *cobra.Command {
	action, past := "disable", "disabled"
	if enabled {
		action, past = "enable", "enabled"
	}
	return &cobra.Command{Use: action + " <set-id>", Short: action + " every Trigger Set member", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		item, err := clientFor(cmd).SetTriggerSetEnabled(ctx, args[0], enabled)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), item)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s Trigger Set %s (%d members)\n", past, item.ID, item.MemberCount)
		return nil
	}}
}

func triggerSetRemove() *cobra.Command {
	return &cobra.Command{Use: "rm <set-id>", Short: "Delete a Trigger Set and all members", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		if err := clientFor(cmd).DeleteTriggerSet(ctx, args[0]); err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), map[string]string{"id": args[0], "status": "deleted"})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted Trigger Set %s\n", args[0])
		return nil
	}}
}
