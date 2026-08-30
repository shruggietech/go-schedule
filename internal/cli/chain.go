package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func newChainCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "chain", Short: "Run tasks after other tasks complete"}
	cmd.AddCommand(chainCreate(), chainList(), chainShow(), chainUpdate(), chainRemove())
	return cmd
}

func chainCreate() *cobra.Command {
	var source, target, outcome string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a task-completion chain",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := validateChainFlags(source, target, outcome); err != nil {
				return err
			}
			ctx, cancel := reqCtx()
			defer cancel()
			chain, err := newClient().CreateChain(ctx, server.ChainCreateRequest{SourceTaskID: source, TargetTaskID: target, OnOutcome: outcome})
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(chain)
			}
			fmt.Fprintf(os.Stdout, "created chain %s: %s (%s) -> %s (%s) on %s\n", chain.ID, chain.SourceTaskName, chain.SourceTaskID, chain.TargetTaskName, chain.TargetTaskID, chain.OnOutcome)
			return nil
		},
	}
	chainFlags(cmd, &source, &target, &outcome)
	return cmd
}

func chainList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List task-completion chains",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			chains, err := newClient().ListChains(ctx)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(chains)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tSOURCE\tTARGET\tON")
			for _, chain := range chains {
				fmt.Fprintf(tw, "%s\t%s (%s)\t%s (%s)\t%s\n", chain.ID, chain.SourceTaskName, chain.SourceTaskID, chain.TargetTaskName, chain.TargetTaskID, chain.OnOutcome)
			}
			return tw.Flush()
		},
	}
}

func chainShow() *cobra.Command {
	return &cobra.Command{
		Use: "show <chain-id>", Short: "Show a task-completion chain", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			chain, err := newClient().GetChain(ctx, args[0])
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(chain)
			}
			fmt.Fprintf(os.Stdout, "%s\nsource: %s (%s)\ntarget: %s (%s)\non: %s\n", chain.ID, chain.SourceTaskName, chain.SourceTaskID, chain.TargetTaskName, chain.TargetTaskID, chain.OnOutcome)
			return nil
		},
	}
}

func chainUpdate() *cobra.Command {
	var source, target, outcome string
	cmd := &cobra.Command{
		Use: "update <chain-id>", Short: "Update a task-completion chain", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := server.ChainUpdateRequest{}
			if cmd.Flags().Changed("source") {
				if source == "" {
					return fmtUsage("--source cannot be empty")
				}
				request.SourceTaskID = &source
			}
			if cmd.Flags().Changed("target") {
				if target == "" {
					return fmtUsage("--target cannot be empty")
				}
				request.TargetTaskID = &target
			}
			if cmd.Flags().Changed("on") {
				if !validChainOutcome(outcome) {
					return fmtUsage("--on must be success, failure, or any")
				}
				request.OnOutcome = &outcome
			}
			if request.SourceTaskID == nil && request.TargetTaskID == nil && request.OnOutcome == nil {
				return fmtUsage("provide at least one of --source, --target, or --on")
			}
			ctx, cancel := reqCtx()
			defer cancel()
			chain, err := newClient().UpdateChain(ctx, args[0], request)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(chain)
			}
			fmt.Fprintf(os.Stdout, "updated chain %s\n", chain.ID)
			return nil
		},
	}
	chainFlags(cmd, &source, &target, &outcome)
	return cmd
}

func chainRemove() *cobra.Command {
	return &cobra.Command{
		Use: "rm <chain-id>", Short: "Delete a task-completion chain", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			if err := newClient().DeleteChain(ctx, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "deleted chain %s\n", args[0])
			return nil
		},
	}
}

func chainFlags(cmd *cobra.Command, source, target, outcome *string) {
	cmd.Flags().StringVar(source, "source", "", "source task ID")
	cmd.Flags().StringVar(target, "target", "", "target task ID")
	cmd.Flags().StringVar(outcome, "on", "", "terminal outcome: success|failure|any")
}

func validateChainFlags(source, target, outcome string) error {
	if source == "" {
		return fmtUsage("--source is required")
	}
	if target == "" {
		return fmtUsage("--target is required")
	}
	if source == target {
		return fmtUsage("--source and --target must identify different tasks")
	}
	if !validChainOutcome(outcome) {
		return fmtUsage("--on must be success, failure, or any")
	}
	return nil
}

func validChainOutcome(outcome string) bool {
	return outcome == "success" || outcome == "failure" || outcome == "any"
}
