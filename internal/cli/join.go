package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func newJoinCmd() *cobra.Command {
	c := &cobra.Command{Use: "join", Short: "Run named parallel branches before one joined target"}
	c.AddCommand(joinCreate(), joinList(), joinShow(), joinUpdate(), joinRemove(), joinInstances(), joinInstance())
	return c
}
func readJoinBranches(cmd *cobra.Command, p string) ([]domain.ParallelJoinBranch, error) {
	if p == "" {
		return nil, fmtUsage("--branches-file is required")
	}
	r := cmd.InOrStdin()
	if p != "-" {
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("open branches file: %w", err)
		}
		defer f.Close()
		r = f
	}
	var rows []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		TaskID string `json:"task_id"`
	}
	var branches []domain.ParallelJoinBranch
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	if err := d.Decode(&rows); err != nil {
		return nil, fmtUsage(fmt.Sprintf("--branches-file must contain a JSON array of id, name and task_id rows: %v", err))
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, fmtUsage("--branches-file must contain exactly one JSON array")
	}
	for _, row := range rows {
		branches = append(branches, domain.ParallelJoinBranch{ID: row.ID, Name: row.Name, TaskID: row.TaskID})
	}
	if len(branches) == 0 {
		return nil, fmtUsage("--branches-file must contain at least one branch")
	}
	ids, names := map[string]bool{}, map[string]bool{}
	for i, b := range branches {
		if strings.TrimSpace(b.Name) == "" || strings.TrimSpace(b.TaskID) == "" {
			return nil, fmtUsage(fmt.Sprintf("branch %d requires name and task_id", i+1))
		}
		if b.ID != "" && ids[b.ID] {
			return nil, fmtUsage(fmt.Sprintf("branch %d repeats id %q", i+1, b.ID))
		}
		if names[b.Name] {
			return nil, fmtUsage(fmt.Sprintf("branch %d repeats name %q", i+1, b.Name))
		}
		ids[b.ID] = true
		names[b.Name] = true
	}
	return branches, nil
}
func joinFlags(c *cobra.Command, name, entry, target, file *string) {
	c.Flags().StringVar(name, "name", "", "parallel block name")
	c.Flags().StringVar(entry, "entry", "", "successful entry task ID")
	c.Flags().StringVar(target, "target", "", "joined target task ID")
	c.Flags().StringVar(file, "branches-file", "", "JSON branch array file, or - for stdin; task references may repeat")
}
func joinCreate() *cobra.Command {
	var name, entry, target, file string
	var routing joinRoutingFlags
	c := &cobra.Command{Use: "create", Short: "Create a parallel join block", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if strings.TrimSpace(name) == "" || entry == "" || target == "" {
			return fmtUsage("--name, --entry and --target are required")
		}
		branches, err := readJoinBranches(cmd, file)
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		request := server.JoinCreateRequest{Name: name, EntryTaskID: entry, TargetTaskID: target, Branches: branches, RecoveryTaskID: routing.recovery, CleanupTaskID: routing.cleanup}
		if cmd.Flags().Changed("recovery-on") {
			conditions, err := joinRecoveryConditions(routing.conditions)
			if err != nil {
				return err
			}
			request.RecoveryOn = &conditions
		}
		v, err := clientFor(cmd).CreateJoin(ctx, request)
		if err != nil {
			return err
		}
		return printAcknowledgement(cmd, v, "created join %s: %s\n", v.ID, v.Name)
	}}
	joinFlags(c, &name, &entry, &target, &file)
	addJoinRoutingFlags(c, &routing)
	return c
}
func joinList() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List parallel join definitions", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		values, err := clientFor(cmd).ListJoins(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), values)
		}
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tENTRY\tBRANCHES\tTARGET\tRECOVERY\tCLEANUP")
		for _, v := range values {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n", v.ID, v.Name, v.EntryTaskID, len(v.Branches), v.TargetTaskID, v.RecoveryTaskID, v.CleanupTaskID)
		}
		return tw.Flush()
	}}
}
func joinShow() *cobra.Command {
	return &cobra.Command{Use: "show <join-id>", Short: "Show a parallel join observation and every branch", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		v, err := clientFor(cmd).GetJoin(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), v)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\nentry: %s (%s)\ntarget: %s (%s)\nupdated: %s\n", v.ID, v.Name, v.EntryTaskName, v.EntryTaskID, v.TargetTaskName, v.TargetTaskID, v.UpdatedAt.Format("2006-01-02T15:04:05.999999999Z07:00")); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "recovery: %s (%s)\nrecovery conditions: %s\ncleanup: %s (%s)\n", v.RecoveryTaskName, v.RecoveryTaskID, joinRecoveryDescription(v.RecoveryOn), v.CleanupTaskName, v.CleanupTaskID); err != nil {
			return err
		}
		if joinObservationTruncated(v) {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Note: some display text is truncated; use --json to identify the affected fields."); err != nil {
				return err
			}
		}
		for _, b := range v.Branches {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "branch %s: %s, task %s (%s)\n", b.ID, b.Name, b.TaskName, b.TaskID); err != nil {
				return err
			}
		}
		return nil
	}}
}
func joinUpdate() *cobra.Command {
	var name, entry, target, file string
	var routing joinRoutingFlags
	c := &cobra.Command{Use: "update <join-id>", Short: "Update future parallel join instances", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		r := server.JoinUpdateRequest{}
		for _, v := range []struct {
			flag  string
			value *string
			dest  **string
		}{{"name", &name, &r.Name}, {"entry", &entry, &r.EntryTaskID}, {"target", &target, &r.TargetTaskID}} {
			if cmd.Flags().Changed(v.flag) {
				if strings.TrimSpace(*v.value) == "" {
					return fmtUsage(fmt.Sprintf("--%s cannot be empty", v.flag))
				}
				*v.dest = v.value
			}
		}
		if cmd.Flags().Changed("branches-file") {
			branches, err := readJoinBranches(cmd, file)
			if err != nil {
				return err
			}
			r.Branches = &branches
		}
		if cmd.Flags().Changed("recovery") {
			r.RecoveryTaskID = &routing.recovery
		}
		if cmd.Flags().Changed("cleanup") {
			r.CleanupTaskID = &routing.cleanup
		}
		if cmd.Flags().Changed("recovery-on") {
			conditions, err := joinRecoveryConditions(routing.conditions)
			if err != nil {
				return err
			}
			r.RecoveryOn = &conditions
		}
		if r.Name == nil && r.EntryTaskID == nil && r.TargetTaskID == nil && r.Branches == nil && r.RecoveryTaskID == nil && r.RecoveryOn == nil && r.CleanupTaskID == nil {
			return fmtUsage("provide --name, --entry, --target, --branches-file, --recovery, --recovery-on or --cleanup")
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		v, err := clientFor(cmd).UpdateJoin(ctx, args[0], r)
		if err != nil {
			return err
		}
		return printAcknowledgement(cmd, v, "updated join %s\n", v.ID)
	}}
	joinFlags(c, &name, &entry, &target, &file)
	addJoinRoutingFlags(c, &routing)
	return c
}
func joinRemove() *cobra.Command {
	return &cobra.Command{Use: "rm <join-id>", Short: "Delete a join definition and resolve outstanding decisions", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		if err := clientFor(cmd).DeleteJoin(ctx, args[0]); err != nil {
			return err
		}
		return printAcknowledgement(cmd, map[string]string{"id": args[0], "status": "deleted"}, "deleted join %s\n", args[0])
	}}
}
func joinInstances() *cobra.Command {
	var id, state, entry string
	var limit, offset int
	c := &cobra.Command{Use: "instances", Short: "Inspect a bounded page of exact join executions", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if limit < 1 || limit > 1000 || offset < 0 {
			return fmtUsage("--limit must be 1 through 1000 and --offset must be nonnegative")
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		v, err := clientFor(cmd).ListJoinInstances(ctx, domain.ParallelJoinInstanceFilter{DefinitionID: id, EntryRunID: entry, State: domain.JoinInstanceState(state), Limit: limit, Offset: offset})
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), v)
		}
		for _, instance := range v {
			if err := printJoinInstance(cmd, instance); err != nil {
				return err
			}
		}
		return nil
	}}
	c.Flags().StringVar(&id, "join", "", "filter by join definition ID")
	c.Flags().StringVar(&state, "state", "", "waiting, blocked, target_pending, recovery_pending, cleanup_pending, completed or resolved")
	c.Flags().StringVar(&entry, "entry-run", "", "filter by exact entry run ID")
	c.Flags().IntVar(&limit, "limit", 100, "maximum instances in this page")
	c.Flags().IntVar(&offset, "offset", 0, "instance page offset")
	return c
}
func joinInstance() *cobra.Command {
	return &cobra.Command{Use: "instance <instance-id>", Short: "Show an exact join instance and every branch", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		v, err := clientFor(cmd).GetJoinInstance(ctx, args[0])
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), v)
		}
		return printJoinInstance(cmd, v)
	}}
}
func printJoinInstance(cmd *cobra.Command, v domain.ParallelJoinInstance) error {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "instance %s: %s (%s)\nstate: %s\nreason: %s\nentry run: %s\nselected route: %s\n", v.ID, v.DefinitionName, v.DefinitionID, v.State, v.Reason, v.EntryRunID, v.SelectedRoute); err != nil {
		return err
	}
	for _, stage := range []struct {
		name, task, delivery, run, reason, cancellation, cancellationError string
		state                                                              domain.DeliveryState
		outcome                                                            domain.RunOutcome
		settlement                                                         domain.JoinSettlement
		started, ended, settled                                            *time.Time
	}{{name: "target", task: v.TargetTaskID, delivery: v.TargetDeliveryID, run: v.TargetRunID, reason: v.TargetReason, state: v.TargetDeliveryState, outcome: v.TargetOutcome, settlement: v.TargetSettlement, started: v.TargetStartedAt, ended: v.TargetEndedAt, settled: v.TargetSettledAt, cancellation: v.TargetCancellationStatus, cancellationError: v.TargetCancellationError},
		{name: "recovery", task: v.RecoveryTaskID, delivery: v.RecoveryDeliveryID, run: v.RecoveryRunID, reason: v.RecoveryReason, state: v.RecoveryDeliveryState, outcome: v.RecoveryOutcome, settlement: v.RecoverySettlement, started: v.RecoveryStartedAt, ended: v.RecoveryEndedAt, settled: v.RecoverySettledAt, cancellation: v.RecoveryCancellationStatus, cancellationError: v.RecoveryCancellationError},
		{name: "cleanup", task: v.CleanupTaskID, delivery: v.CleanupDeliveryID, run: v.CleanupRunID, reason: v.CleanupReason, state: v.CleanupDeliveryState, outcome: v.CleanupOutcome, settlement: v.CleanupSettlement, started: v.CleanupStartedAt, ended: v.CleanupEndedAt, settled: v.CleanupSettledAt, cancellation: v.CleanupCancellationStatus, cancellationError: v.CleanupCancellationError}} {
		if stage.task == "" && stage.name != "target" {
			continue
		}
		nativeState := "waiting"
		if stage.started != nil && stage.ended == nil {
			nativeState = "running"
		} else if stage.settlement != "" {
			nativeState = string(stage.settlement)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s task: %s\n%s admission: %s %s\n%s run: %s, outcome: %s, native state: %s\n%s started: %s, ended: %s\n%s settlement: %s at %s, reason: %s\n", stage.name, stage.task, stage.name, stage.delivery, stage.state, stage.name, stage.run, stage.outcome, nativeState, stage.name, joinObservationTime(stage.started), joinObservationTime(stage.ended), stage.name, stage.settlement, joinObservationTime(stage.settled), stage.reason); err != nil {
			return err
		}
		if stage.cancellation != "" {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s cancellation request: %s, detail: %s\n", stage.name, stage.cancellation, stage.cancellationError); err != nil {
				return err
			}
		}
	}
	for _, branch := range v.Branches {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "branch %s (%s), invocation %s, task %s, delivery %s %s, run %s, outcome %s: %s\n  settlement: %s at %s, started: %s, ended: %s\n", branch.Name, branch.BranchID, branch.ID, branch.TaskID, branch.DeliveryID, branch.DeliveryState, branch.RunID, branch.Outcome, branch.Reason, branch.Settlement, joinObservationTime(branch.SettledAt), joinObservationTime(branch.StartedAt), joinObservationTime(branch.EndedAt)); err != nil {
			return err
		}
		if branch.CancellationStatus != "" {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  cancellation request: %s, detail: %s\n", branch.CancellationStatus, branch.CancellationError); err != nil {
				return err
			}
		}
	}
	return nil
}

func joinObservationTime(value *time.Time) string {
	if value == nil {
		return "not reported"
	}
	return value.Format(time.RFC3339Nano)
}

func joinObservationTruncated(v domain.ParallelJoin) bool {
	if len(v.TruncatedFields) > 0 {
		return true
	}
	for _, branch := range v.Branches {
		if len(branch.TruncatedFields) > 0 {
			return true
		}
	}
	return false
}

// joinRoutingFlags extends the existing join family with optional continuation intent.
type joinRoutingFlags struct {
	recovery, cleanup string
	conditions        []string
}

func addJoinRoutingFlags(c *cobra.Command, values *joinRoutingFlags) {
	c.Flags().StringVar(&values.recovery, "recovery", "", "recovery task ID; on update an empty value clears recovery")
	c.Flags().StringSliceVar(&values.conditions, "recovery-on", nil, "comma-separated failure,canceled,skipped,not_started; omitted defaults to failure, none selects no causes")
	c.Flags().StringVar(&values.cleanup, "cleanup", "", "cleanup task ID; runs after branches and selected follow-up settle; empty clears on update")
}
func joinRecoveryConditions(values []string) ([]domain.JoinRecoveryCondition, error) {
	out := []domain.JoinRecoveryCondition{}
	seen := map[string]bool{}
	for _, value := range values {
		if value == "none" {
			if len(values) != 1 {
				return nil, fmtUsage("--recovery-on none cannot be combined with other conditions")
			}
			return out, nil
		}
		switch domain.JoinRecoveryCondition(value) {
		case domain.JoinRecoveryOnFailure, domain.JoinRecoveryOnCanceled, domain.JoinRecoveryOnSkipped, domain.JoinRecoveryOnNotStarted:
		default:
			return nil, fmtUsage("--recovery-on accepts failure,canceled,skipped,not_started or none")
		}
		if seen[value] {
			return nil, fmtUsage("--recovery-on must not repeat a condition")
		}
		seen[value] = true
		out = append(out, domain.JoinRecoveryCondition(value))
	}
	return out, nil
}
func joinRecoveryDescription(values []domain.JoinRecoveryCondition) string {
	if len(values) == 0 {
		return "none"
	}
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = string(v)
	}
	return strings.Join(out, ",")
}
