package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// missingDateUsage is shared by add and edit so the two cannot drift. The values
// are underscored to match --overlap and --catchup rather than hyphenated to
// match the flag name; consistency across the policy flags wins.
const missingDateUsage = "what to do in a period with no matching date: skip|last_valid|next_valid"

func newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "task", Short: "Create and manage tasks"}
	cmd.AddCommand(taskAdd(), taskList(), taskShow(), taskEdit(), taskEnable(), taskDisable(), taskRm(), taskRunNow())
	return cmd
}

// groupIntent maps the --group flag onto the API's three-way group membership
// intent. Omitting the flag leaves membership alone (nil), which is the
// long-standing behavior of `task edit`; passing it explicitly assigns the given
// group, or removes the task from its group when the value is empty.
func groupIntent(cmd *cobra.Command, group string) *string {
	if !cmd.Flags().Changed("group") {
		return nil
	}
	return &group
}

func taskEdit() *cobra.Command {
	var command, cwd, group, tz, sched, at, overlap, catchup, missingDate string
	var args, env []string
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Modify a task (only provided fields change)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			if sched != "" && at != "" {
				return fmt.Errorf("%w: provide at most one of --schedule or --at", errUsage)
			}
			envMap, err := parseEnv(env)
			if err != nil {
				return fmt.Errorf("%w: %v", errUsage, err)
			}
			req := server.TaskUpdateRequest{
				Command: command, Args: args, WorkingDir: cwd, Env: envMap,
				Timezone: tz, Schedule: sched, OverlapPolicy: overlap, CatchupPolicy: catchup,
				MissingDatePolicy: missingDate,
			}
			req.GroupID = groupIntent(cmd, group)
			if at != "" {
				ts, err := time.Parse(time.RFC3339, at)
				if err != nil {
					return fmt.Errorf("%w: --at must be RFC3339", errUsage)
				}
				req.At = &ts
			}
			ctx, cancel := reqCtx()
			defer cancel()
			resp, err := newClient().UpdateTask(ctx, a[0], req)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(resp)
			}
			fmt.Fprintf(os.Stdout, "updated task %s\nschedule: %s\n", resp.Task.ID, resp.Schedule.HumanSummary)
			printNextRuns(resp.NextRuns)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&command, "command", "", "program or script to run")
	f.StringArrayVar(&args, "arg", nil, "argument (repeatable; replaces existing)")
	f.StringVar(&cwd, "cwd", "", "working directory")
	f.StringArrayVar(&env, "env", nil, "environment variable KEY=VALUE (repeatable; replaces existing)")
	f.StringVar(&group, "group", "", `group ID (pass "" to remove the task from its group)`)
	f.StringVar(&tz, "tz", "", "IANA timezone")
	f.StringVar(&sched, "schedule", "", "new human-readable recurrence")
	f.StringVar(&at, "at", "", "new one-off run time (RFC3339)")
	f.StringVar(&overlap, "overlap", "", "overlap policy: queue_one|skip|allow_concurrent")
	f.StringVar(&catchup, "catchup", "", "catch-up policy: one|none")
	f.StringVar(&missingDate, "missing-date", "", missingDateUsage)
	return cmd
}

func taskAdd() *cobra.Command {
	var (
		command     string
		args        []string
		cwd         string
		env         []string
		group       string
		tz          string
		sched       string
		at          string
		overlap     string
		catchup     string
		missingDate string
	)
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a task (recurring via --schedule or one-off via --at)",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			if command == "" {
				return fmt.Errorf("%w: --command is required", errUsage)
			}
			if (sched == "") == (at == "") {
				return fmt.Errorf("%w: provide exactly one of --schedule or --at", errUsage)
			}
			envMap, err := parseEnv(env)
			if err != nil {
				return fmt.Errorf("%w: %v", errUsage, err)
			}
			req := server.TaskCreateRequest{
				Name: a[0], Command: command, Args: args, WorkingDir: cwd, Env: envMap,
				GroupID: group, Timezone: tz, Schedule: sched, OverlapPolicy: overlap, CatchupPolicy: catchup,
				MissingDatePolicy: missingDate,
			}
			if at != "" {
				ts, err := time.Parse(time.RFC3339, at)
				if err != nil {
					return fmt.Errorf("%w: --at must be RFC3339 (e.g. 2026-08-04T09:00:00Z)", errUsage)
				}
				req.At = &ts
			}
			ctx, cancel := reqCtx()
			defer cancel()
			resp, err := newClient().CreateTask(ctx, req)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(resp)
			}
			fmt.Fprintf(os.Stdout, "created task %s (%s)\nschedule: %s\n", resp.Task.ID, resp.Task.Name, resp.Schedule.HumanSummary)
			printNextRuns(resp.NextRuns)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&command, "command", "", "program or script to run (required)")
	f.StringArrayVar(&args, "arg", nil, "argument to the command (repeatable)")
	f.StringVar(&cwd, "cwd", "", "working directory")
	f.StringArrayVar(&env, "env", nil, "environment variable KEY=VALUE (repeatable)")
	f.StringVar(&group, "group", "", "group ID")
	f.StringVar(&tz, "tz", "", "IANA timezone (default: system local)")
	f.StringVar(&sched, "schedule", "", `human-readable recurrence, e.g. "every 15 minutes", "3rd wednesday monthly at 14:00"`)
	f.StringVar(&at, "at", "", "one-off run time (RFC3339)")
	f.StringVar(&overlap, "overlap", "", "overlap policy: queue_one|skip|allow_concurrent")
	f.StringVar(&catchup, "catchup", "", "catch-up policy: one|none")
	f.StringVar(&missingDate, "missing-date", "", missingDateUsage)
	return cmd
}

func taskList() *cobra.Command {
	var group, state string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			tasks, err := newClient().ListTasks(ctx, group, state)
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(tasks)
			}
			tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tNAME\tSTATE\tENABLED\tTZ")
			for _, t := range tasks {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%t\t%s\n", t.ID, t.Name, t.State, t.Enabled, t.Timezone)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().StringVar(&group, "group", "", "filter by group ID")
	cmd.Flags().StringVar(&state, "state", "", "filter by state (active|completed|disabled)")
	return cmd
}

func taskShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show task detail and upcoming runs",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			resp, err := newClient().GetTask(ctx, a[0])
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(resp)
			}
			t := resp.Task
			fmt.Fprintf(os.Stdout, "%s  %s\nstate: %s  enabled: %t  tz: %s\ncommand: %s %s\nschedule: %s\noverlap: %s  catch-up: %s  missing dates: %s\n",
				t.ID, t.Name, t.State, t.Enabled, t.Timezone, t.Command, strings.Join(t.Args, " "),
				resp.Schedule.HumanSummary, t.OverlapPolicy, t.CatchupPolicy, t.MissingDatePolicy)
			printNextRuns(resp.NextRuns)
			return nil
		},
	}
}

func taskEnable() *cobra.Command {
	return &cobra.Command{
		Use: "enable <id>", Short: "Enable a task", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error { return toggle(a[0], true) },
	}
}

func taskDisable() *cobra.Command {
	return &cobra.Command{
		Use: "disable <id>", Short: "Disable a task", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error { return toggle(a[0], false) },
	}
}

func toggle(id string, enabled bool) error {
	ctx, cancel := reqCtx()
	defer cancel()
	if err := newClient().SetTaskEnabled(ctx, id, enabled); err != nil {
		return err
	}
	state := "disabled"
	if enabled {
		state = "enabled"
	}
	fmt.Fprintf(os.Stdout, "task %s %s\n", id, state)
	return nil
}

func taskRm() *cobra.Command {
	return &cobra.Command{
		Use: "rm <id>", Short: "Delete a task", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			if err := newClient().DeleteTask(ctx, a[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "deleted task %s\n", a[0])
			return nil
		},
	}
}

func taskRunNow() *cobra.Command {
	return &cobra.Command{
		Use: "run-now <id>", Short: "Trigger an immediate run", Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, a []string) error {
			ctx, cancel := reqCtx()
			defer cancel()
			if err := newClient().RunNow(ctx, a[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "triggered run for task %s\n", a[0])
			return nil
		},
	}
}

func parseEnv(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("invalid --env %q (want KEY=VALUE)", p)
		}
		m[k] = v
	}
	return m, nil
}

func printNextRuns(runs []time.Time) {
	if len(runs) == 0 {
		return
	}
	fmt.Fprintln(os.Stdout, "next runs:")
	for _, r := range runs {
		fmt.Fprintf(os.Stdout, "  %s\n", r.Format(time.RFC3339))
	}
}
