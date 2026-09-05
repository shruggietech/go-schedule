package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func newWatcherCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "watcher", Aliases: []string{"watch"}, Short: "Manage filesystem watchers"}
	cmd.AddCommand(newWatcherCreateCmd(), newWatcherListCmd(), newWatcherShowCmd(), newWatcherUpdateCmd(), newWatcherEnabledCmd(true), newWatcherEnabledCmd(false), newWatcherDeleteCmd())
	return cmd
}

func newWatcherCreateCmd() *cobra.Command {
	var req server.FilesystemWatcherCreateRequest
	var enabled bool
	cmd := &cobra.Command{Use: "create", Short: "Create a filesystem watcher", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		req.Enabled = &enabled
		ctx, cancel := reqCtx()
		defer cancel()
		result, err := newClient().CreateFilesystemWatcher(ctx, req)
		if err != nil {
			return err
		}
		return printWatcher(result)
	}}
	cmd.Flags().StringVar(&req.Name, "name", "", "watcher name")
	cmd.Flags().Var((*watcherKindValue)(&req.Kind), "kind", "file or directory")
	cmd.Flags().StringVar(&req.Path, "path", "", "absolute or relative path")
	cmd.Flags().StringVar(&req.Pattern, "pattern", "", "directory basename glob")
	cmd.Flags().BoolVar(&req.Recursive, "recursive", false, "include descendant directories")
	cmd.Flags().StringVar(&req.Debounce, "debounce", "", "quiet interval such as 250ms")
	cmd.Flags().StringVar(&req.Stability, "stability", "", "file stability interval such as 500ms")
	cmd.Flags().StringVar(&req.TargetTaskID, "task", "", "target task ID")
	cmd.Flags().BoolVar(&enabled, "enabled", true, "enable on creation")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("kind")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("task")
	return cmd
}

type watcherKindValue domain.WatcherKind

func (value *watcherKindValue) String() string { return string(*value) }
func (value *watcherKindValue) Type() string   { return "watcher-kind" }
func (value *watcherKindValue) Set(raw string) error {
	kind := domain.WatcherKind(raw)
	if kind != domain.WatcherFile && kind != domain.WatcherDirectory {
		return fmtUsage("kind must be file or directory")
	}
	*value = watcherKindValue(kind)
	return nil
}

func newWatcherListCmd() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List filesystem watchers", Args: cobra.NoArgs, RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		items, err := newClient().ListFilesystemWatchers(ctx)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(items)
		}
		for _, item := range items {
			fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\t%s\n", item.ID, item.Name, item.Kind, item.Readiness, item.Path)
		}
		return nil
	}}
}

func newWatcherShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <watcher-id>", Short: "Show a filesystem watcher", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		item, err := newClient().GetFilesystemWatcher(ctx, args[0])
		if err != nil {
			return err
		}
		return printWatcher(item)
	}}
}

func newWatcherUpdateCmd() *cobra.Command {
	var name, path, pattern, debounce, stability, task string
	var kind domain.WatcherKind
	var recursive bool
	cmd := &cobra.Command{Use: "update <watcher-id>", Short: "Update a filesystem watcher", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		var req server.FilesystemWatcherUpdateRequest
		if cmd.Flags().Changed("name") {
			req.Name = &name
		}
		if cmd.Flags().Changed("kind") {
			req.Kind = &kind
		}
		if cmd.Flags().Changed("path") {
			req.Path = &path
		}
		if cmd.Flags().Changed("pattern") {
			req.Pattern = &pattern
		}
		if cmd.Flags().Changed("recursive") {
			req.Recursive = &recursive
		}
		if cmd.Flags().Changed("debounce") {
			req.Debounce = &debounce
		}
		if cmd.Flags().Changed("stability") {
			req.Stability = &stability
		}
		if cmd.Flags().Changed("task") {
			req.TargetTaskID = &task
		}
		ctx, cancel := reqCtx()
		defer cancel()
		item, err := newClient().UpdateFilesystemWatcher(ctx, args[0], req)
		if err != nil {
			return err
		}
		return printWatcher(item)
	}}
	cmd.Flags().StringVar(&name, "name", "", "watcher name")
	cmd.Flags().Var((*watcherKindValue)(&kind), "kind", "file or directory")
	cmd.Flags().StringVar(&path, "path", "", "path")
	cmd.Flags().StringVar(&pattern, "pattern", "", "directory basename glob")
	cmd.Flags().BoolVar(&recursive, "recursive", false, "include descendant directories")
	cmd.Flags().StringVar(&debounce, "debounce", "", "quiet interval")
	cmd.Flags().StringVar(&stability, "stability", "", "stability interval")
	cmd.Flags().StringVar(&task, "task", "", "target task ID")
	return cmd
}

func newWatcherEnabledCmd(enabled bool) *cobra.Command {
	action := "disable"
	if enabled {
		action = "enable"
	}
	return &cobra.Command{Use: action + " <watcher-id>", Short: action + " a filesystem watcher", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		item, err := newClient().SetFilesystemWatcherEnabled(ctx, args[0], enabled)
		if err != nil {
			return err
		}
		return printWatcher(item)
	}}
}

func newWatcherDeleteCmd() *cobra.Command {
	return &cobra.Command{Use: "rm <watcher-id>", Short: "Delete a filesystem watcher", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		return newClient().DeleteFilesystemWatcher(ctx, args[0])
	}}
}

func printWatcher(item server.FilesystemWatcherResponse) error {
	if jsonOut {
		return printJSON(item)
	}
	return printWatcherTo(os.Stdout, item)
}

func printWatcherTo(output io.Writer, item server.FilesystemWatcherResponse) error {
	_, err := fmt.Fprintf(output, "ID: %s\nName: %s\nKind: %s\nPath: %s\nPattern: %s\nRecursive: %t\nDebounce: %s\nStability: %s\nTarget: %s (%s)\nEnabled: %t\nHealth: %s\nHealth reason: %s\nReadiness: %s\nReadiness reason: %s\n", item.ID, item.Name, item.Kind, item.Path, item.Pattern, item.Recursive, item.Debounce, item.Stability, item.TargetTaskName, item.TargetTaskID, item.Enabled, item.Health.State, item.Health.Reason, item.Readiness, item.Reason)
	return err
}
