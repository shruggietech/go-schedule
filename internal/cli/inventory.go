package cli

import (
	"encoding/json"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
)

type inventoryFlag struct {
	Name       string `json:"name"`
	Shorthand  string `json:"shorthand"`
	Type       string `json:"type"`
	Default    string `json:"default"`
	Usage      string `json:"usage"`
	Scope      string `json:"scope"`
	Hidden     bool   `json:"hidden"`
	Deprecated string `json:"deprecated"`
	Required   bool   `json:"required"`
}
type inventoryCommand struct {
	Path     string          `json:"path"`
	Use      string          `json:"use"`
	Short    string          `json:"short"`
	Long     string          `json:"long"`
	Aliases  []string        `json:"aliases"`
	Hidden   bool            `json:"hidden"`
	Runnable bool            `json:"runnable"`
	Flags    []inventoryFlag `json:"flags"`
}

// CommandInventoryJSON describes the actual initialized tree without host I/O.
// Only Cobra's generated metadata completion initializes hidden callbacks; product handlers and host discovery never run.
func CommandInventoryJSON() ([]byte, error) {
	root := newRootWithDependencies(Dependencies{Client: client.New("inventory-only")})
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{cobra.ShellCompRequestCmd, ""})
	if err := root.Execute(); err != nil {
		return nil, err
	}
	commands := []inventoryCommand{}
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		cmd.InitDefaultHelpFlag()
		cmd.InheritedFlags()
		record := inventoryCommand{Path: strings.TrimSpace(strings.TrimPrefix(cmd.CommandPath(), root.Name())), Use: cmd.Use, Short: cmd.Short, Long: cmd.Long, Aliases: append([]string{}, cmd.Aliases...), Hidden: cmd.Hidden, Runnable: cmd.Runnable(), Flags: []inventoryFlag{}}
		sort.Strings(record.Aliases)
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			scope := "local"
			if cmd.InheritedFlags().Lookup(flag.Name) != nil {
				scope = "inherited"
			}
			if cmd.PersistentFlags().Lookup(flag.Name) != nil {
				scope = "persistent"
			}
			required := false
			for _, value := range flag.Annotations[cobra.BashCompOneRequiredFlag] {
				required = required || value == "true"
			}
			record.Flags = append(record.Flags, inventoryFlag{Name: flag.Name, Shorthand: flag.Shorthand, Type: flag.Value.Type(), Default: flag.DefValue, Usage: flag.Usage, Scope: scope, Hidden: flag.Hidden, Deprecated: flag.Deprecated, Required: required})
		})
		commands = append(commands, record)
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)
	sort.Slice(commands, func(i, j int) bool { return commands[i].Path < commands[j].Path })
	return json.MarshalIndent(commands, "", "  ")
}
