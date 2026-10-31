package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

// taskInputFlags carries execution-input intent separately from omitted values.
type taskInputFlags struct {
	stdin, stdinFile, runAs, syntax string
	clearStdin, clearRunAs          bool
	clearArgs, clearCwd, clearEnv   bool
	disabled                        bool
}

func (f *taskInputFlags) bind(cmd *cobra.Command, edit bool) {
	flags := cmd.Flags()
	flags.StringVar(&f.stdin, "stdin", "", "text sent to the task's standard input (empty clears it)")
	flags.StringVar(&f.stdinFile, "stdin-file", "", "UTF-8 task input file; - reads this command's standard input")
	flags.StringVar(&f.runAs, "run-as", "", "execution user supported by the target operating system")
	flags.StringVar(&f.syntax, "syntax", "auto", "schedule interpretation: auto|human|cron (requires --schedule)")
	if edit {
		flags.BoolVar(&f.clearStdin, "clear-stdin", false, "remove the task's standard input")
		flags.BoolVar(&f.clearRunAs, "clear-run-as", false, "restore the daemon's default execution identity")
		flags.BoolVar(&f.clearArgs, "clear-args", false, "remove all command arguments")
		flags.BoolVar(&f.clearCwd, "clear-cwd", false, "restore the default working directory")
		flags.BoolVar(&f.clearEnv, "clear-env", false, "remove task environment overrides; keep the base process environment")
	} else {
		flags.BoolVar(&f.disabled, "disabled", false, "create disabled in one operation, without an enabled interval")
	}
}

func (f taskInputFlags) validate(cmd *cobra.Command) error {
	if cmd.Flags().Changed("stdin") && cmd.Flags().Changed("stdin-file") {
		return fmtUsage("use only one of --stdin and --stdin-file")
	}
	for _, pair := range []struct {
		clear bool
		flag  string
		name  string
	}{
		{f.clearStdin, "stdin", "clear-stdin"},
		{f.clearStdin, "stdin-file", "clear-stdin"},
		{f.clearRunAs, "run-as", "clear-run-as"},
		{f.clearArgs, "arg", "clear-args"},
		{f.clearCwd, "cwd", "clear-cwd"},
		{f.clearEnv, "env", "clear-env"},
	} {
		if pair.clear && cmd.Flags().Changed(pair.flag) {
			return fmtUsage(fmt.Sprintf("use only one of --%s and --%s", pair.flag, pair.name))
		}
	}
	if cmd.Flags().Changed("syntax") {
		value, err := cmd.Flags().GetString("schedule")
		if err != nil {
			return err
		}
		if strings.TrimSpace(value) == "" {
			return fmtUsage("--syntax requires a nonempty --schedule")
		}
	}
	_, err := scheduleSyntax(f.syntax)
	return err
}

func scheduleSyntax(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "auto":
		return "", nil
	case "human":
		return "human", nil
	case "cron":
		return "cron", nil
	default:
		return "", fmtUsage("--syntax must be auto, human, or cron")
	}
}

func (f taskInputFlags) input(cmd *cobra.Command) (*string, error) {
	if f.clearStdin {
		value := ""
		return &value, nil
	}
	if cmd.Flags().Changed("stdin-file") {
		var data []byte
		var err error
		if f.stdinFile == "-" {
			data, err = io.ReadAll(cmd.InOrStdin())
		} else {
			data, err = os.ReadFile(f.stdinFile)
		}
		if err != nil {
			return nil, fmt.Errorf("read --stdin-file: %w", err)
		}
		if !utf8.Valid(data) {
			return nil, fmtUsage("--stdin-file must contain UTF-8 text; save the input as UTF-8 and try again")
		}
		value := string(data)
		return &value, nil
	}
	if cmd.Flags().Changed("stdin") {
		return &f.stdin, nil
	}
	return nil, nil
}

func (f taskInputFlags) create(cmd *cobra.Command, request *server.TaskCreateRequest) error {
	if err := f.validate(cmd); err != nil {
		return err
	}
	input, err := f.input(cmd)
	if err != nil {
		return err
	}
	if input != nil {
		request.Stdin = *input
	}
	request.RunAs = f.runAs
	request.ScheduleSyntax, err = scheduleSyntax(f.syntax)
	if cmd.Flags().Changed("disabled") {
		enabled := !f.disabled
		request.Enabled = &enabled
	}
	return err
}

func (f taskInputFlags) update(cmd *cobra.Command, request *server.TaskUpdateRequest) error {
	if err := f.validate(cmd); err != nil {
		return err
	}
	var err error
	request.Stdin, err = f.input(cmd)
	if err != nil {
		return err
	}
	request.RunAs = f.runAs
	request.ClearRunAs = f.clearRunAs || cmd.Flags().Changed("run-as") && f.runAs == ""
	request.ClearWorkingDir = f.clearCwd || cmd.Flags().Changed("cwd") && request.WorkingDir == ""
	request.ClearEnv = f.clearEnv
	if f.clearArgs {
		request.Args = []string{}
	}
	request.ScheduleSyntax, err = scheduleSyntax(f.syntax)
	return err
}
