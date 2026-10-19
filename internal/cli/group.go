package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/task"
)

func newGroupCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "group", Short: "Organize tasks into nested groups"}
	cmd.AddCommand(groupAdd(), groupList(), groupEdit(), groupEnable(), groupDisable(), groupRm())
	return cmd
}

func groupAdd() *cobra.Command {
	var parent string
	var disabled bool
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a group (optionally under --parent)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			request := server.GroupCreateRequest{Name: a[0], ParentID: parent}
			if cmd.Flags().Changed("disabled") {
				enabled := !disabled
				request.Enabled = &enabled
			}
			g, err := clientFor(cmd).CreateGroup(ctx, request)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), g)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created group %s (%s)\n", g.ID, g.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&parent, "parent", "", "parent group ID")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "create disabled in one operation, without an enabled interval")
	return cmd
}

func groupEdit() *cobra.Command {
	var name, parent string
	cmd := &cobra.Command{
		Use: "edit <id>", Short: "Rename a group or change its parent", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("parent") {
				return fmtUsage("provide --name or --parent; use --parent= to move the group to the root")
			}
			request := server.GroupUpdateRequest{}
			if cmd.Flags().Changed("name") {
				if strings.TrimSpace(name) == "" {
					return fmtUsage("--name must not be empty; provide the new group name")
				}
				request.Name = strings.TrimSpace(name)
			}
			if cmd.Flags().Changed("parent") {
				request.Parent = &parent
			}
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			group, err := clientFor(cmd).UpdateGroup(ctx, args[0], request)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), group)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "updated group %s (%s), parent=%q\n", group.ID, group.Name, group.ParentID)
			return err
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "new nonempty group name")
	cmd.Flags().StringVar(&parent, "parent", "", "new parent group ID; empty moves the group to the root")
	return cmd
}

func groupList() *cobra.Command {
	var tree bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			cl := clientFor(cmd)
			if tree {
				forest, err := cl.GroupTree(ctx)
				if err != nil {
					return err
				}
				if jsonOutput(cmd) {
					return printJSONTo(cmd.OutOrStdout(), forest)
				}
				for _, n := range forest {
					printTree(cmd, n, 0)
				}
				return nil
			}
			groups, err := cl.ListGroups(ctx)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return printJSONTo(cmd.OutOrStdout(), groups)
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tNAME\tPARENT\tENABLED")
			for _, g := range groups {
				parent := g.ParentID
				if parent == "" {
					parent = "-"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%t\n", g.ID, g.Name, parent, g.Enabled)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().BoolVar(&tree, "tree", false, "show as a hierarchy")
	return cmd
}

func printTree(cmd *cobra.Command, n *task.TreeNode, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	state := ""
	if !n.Group.Enabled {
		state = " (disabled)"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s- %s [%s]%s\n", indent, n.Group.Name, n.Group.ID, state)
	for _, c := range n.Children {
		printTree(cmd, c, depth+1)
	}
}

func groupEnable() *cobra.Command {
	return &cobra.Command{
		Use: "enable <id>", Short: "Enable a group (and its subtree)", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error { return groupToggle(cmd, a[0], true) },
	}
}

func groupDisable() *cobra.Command {
	return &cobra.Command{
		Use: "disable <id>", Short: "Disable a group (cascades to its subtree)", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error { return groupToggle(cmd, a[0], false) },
	}
}

func groupToggle(cmd *cobra.Command, id string, enabled bool) error {
	ctx, cancel := reqCtx(cmd)
	defer cancel()
	if err := clientFor(cmd).SetGroupEnabled(ctx, id, enabled); err != nil {
		return err
	}
	state := "disabled"
	if enabled {
		state = "enabled"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "group %s %s\n", id, state)
	return nil
}

func groupRm() *cobra.Command {
	return &cobra.Command{
		Use: "rm <id>", Short: "Delete a group (children cascade; tasks are ungrouped)", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, a []string) error {
			ctx, cancel := reqCtx(cmd)
			defer cancel()
			if err := clientFor(cmd).DeleteGroup(ctx, a[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted group %s\n", a[0])
			return nil
		},
	}
}
