package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// executeInvocation classifies command lookup separately from operational errors.
// Find does not parse flags or invoke callbacks, so execution still parses once.
func executeInvocation(ctx context.Context, root *cobra.Command, args []string) int {
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	if len(args) > 0 && (args[0] == cobra.ShellCompRequestCmd || args[0] == cobra.ShellCompNoDescRequestCmd) {
		// Cobra installs these hidden protocol commands only inside ExecuteC.
		// Their raw completion words must bypass ordinary command lookup.
		if len(args) < 2 {
			return handleExecuteError(root.ErrOrStderr(), fmtUsage("completion protocol requires a word to complete (an empty word is allowed)"))
		}
		root.SetArgs(args)
		return handleExecuteError(root.ErrOrStderr(), root.ExecuteContext(ctx))
	}
	cmd, remaining, err := root.Find(args)
	if err != nil {
		return handleExecuteError(root.ErrOrStderr(), usageCause(err))
	}
	if !cmd.Runnable() {
		// Cobra returns help before validating namespace operands. Parse those
		// flags once, reject unknown subcommands, then execute only the path so
		// help/version handling cannot append repeatable flags a second time.
		cmd.InitDefaultHelpFlag()
		cmd.InitDefaultVersionFlag()
		if err := cmd.ParseFlags(remaining); err != nil {
			return handleExecuteError(root.ErrOrStderr(), usageCause(err))
		}
		if cmd.Flags().Lookup("help").Value.String() != "true" {
			if err := cobra.NoArgs(cmd, cmd.Flags().Args()); err != nil {
				return handleExecuteError(root.ErrOrStderr(), usageCause(err))
			}
		}
		args = strings.Fields(strings.TrimPrefix(cmd.CommandPath(), root.Name()))
	}
	root.SetArgs(args)
	err = root.ExecuteContext(ctx)
	if err == nil && stateFor(root).helpError != nil {
		err = fmt.Errorf("write command help: %w", stateFor(root).helpError)
	}
	return handleExecuteError(root.ErrOrStderr(), err)
}

func usageCause(err error) error {
	if errors.Is(err, errUsage) {
		return err
	}
	return fmt.Errorf("%w: %w", errUsage, err)
}

func configureCommandContracts(root *cobra.Command) {
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	configurePowerShellCompletionHelp(root)
	configureLocalHelp(root)
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageCause(err) })
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd.Runnable() {
			validate := cmd.Args
			if validate == nil {
				validate = cobra.NoArgs
			}
			cmd.Args = func(cmd *cobra.Command, args []string) error {
				if err := validate(cmd, args); err != nil {
					return usageCause(err)
				}
				return nil
			}
			if run := cmd.RunE; run != nil {
				cmd.RunE = func(cmd *cobra.Command, args []string) error {
					return contextualCommandError(cmd, args, run(cmd, args))
				}
			}
		}
		// Namespace Args stays nil: Cobra Find then rejects unknown subcommands.
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)
}

func configureLocalHelp(root *cobra.Command) {
	standardHelp := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		output := cmd.OutOrStdout()
		tracked := &trackingWriter{writer: output}
		cmd.SetOut(tracked)
		defer cmd.SetOut(output)
		standardHelp(cmd, args)
		stateFor(cmd).helpError = tracked.err
	})
	for _, help := range root.Commands() {
		if help.Name() != "help" {
			continue
		}
		help.Args = cobra.ArbitraryArgs
		help.Run = nil
		help.RunE = func(_ *cobra.Command, args []string) error {
			topic, remaining, err := root.Find(args)
			if err != nil {
				return usageCause(err)
			}
			if len(remaining) != 0 {
				return fmtUsage("unknown help topic; use gosched help to list commands")
			}
			return topic.Help()
		}
	}
}

func configurePowerShellCompletionHelp(root *cobra.Command) {
	for _, command := range root.Commands() {
		if command.Name() != "completion" {
			continue
		}
		for _, shell := range command.Commands() {
			if shell.Name() == "powershell" {
				shell.Long = fmt.Sprintf(`Generate the autocompletion script for PowerShell.

To load completions in your current PowerShell session:

	%s completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the following invocation as a line in your PowerShell profile ($PROFILE):

	%s completion powershell | Out-String | Invoke-Expression
`, root.Name(), root.Name())
			}
		}
	}
}
