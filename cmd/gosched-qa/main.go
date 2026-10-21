// gosched-qa captures isolated maintainer evaluations and validates agent reviews.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clievaluate"
	"github.com/shruggietech/go-schedule/internal/platform"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "gosched-qa:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out, diagnostic io.Writer) error {
	if len(args) == 0 {
		return errors.New("choose build, inventory, capture, resume, finalize, validate or render; use COMMAND --help for its inputs")
	}
	command := args[0]
	flags := flag.NewFlagSet("gosched-qa "+command, flag.ContinueOnError)
	flags.SetOutput(diagnostic)
	if command == "fixture-child" {
		mode := flags.String("mode", "success", "owned fixture behavior")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		switch *mode {
		case "tree":
			executable, err := os.Executable()
			if err != nil {
				return err
			}
			directory, err := os.Getwd()
			if err != nil {
				return err
			}
			child := exec.Command(executable, "fixture-child", "--mode=hold")
			child.Dir = directory
			child.Env = os.Environ()
			child.Stdin = nil
			child.Stdout = out
			child.Stderr = diagnostic
			platform.HideConsole(child)
			if err := child.Start(); err != nil {
				return err
			}
			defer func() { _ = child.Process.Kill(); _ = child.Wait() }()
			if _, err := fmt.Fprintf(out, "owned descendant %d\n", child.Process.Pid); err != nil {
				return err
			}
			<-ctx.Done()
			return ctx.Err()
		case "hold":
			<-ctx.Done()
			return ctx.Err()
		case "failure":
			fmt.Fprintln(diagnostic, "fixture deliberately failed with exit 7")
			os.Exit(7)
		case "output":
			_, err := io.WriteString(out, strings.Repeat("fixture output\n", 10000))
			return err
		case "wait":
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 * time.Second):
			}
		case "success":
		default:
			return errors.New("unknown owned fixture mode")
		}
		input, err := io.ReadAll(io.LimitReader(os.Stdin, 4096))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "fixture completed args=%q stdin=%q QA_VALUE=%q\n", flags.Args(), input, os.Getenv("QA_VALUE"))
		return err
	}
	if command == "inventory" {
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if flags.NArg() != 0 {
			return errors.New("inventory accepts no positional input")
		}
		manifest, _, err := clievaluate.BuildManifest()
		if err != nil {
			return err
		}
		data, err := clievaluate.Encode(manifest)
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		return err
	}
	if command == "worker" {
		path := flags.String("config", "", "owned coordinator configuration")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		var config clievaluate.WorkerConfig
		err := clievaluate.ReadInput(*path, &config)
		if err != nil {
			return err
		}
		workerCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Budget.TotalTimeoutMS)*time.Millisecond)
		defer cancel()
		go func() { _, _ = io.Copy(io.Discard, os.Stdin); cancel() }()
		return clievaluate.RunWorker(workerCtx, config)
	}
	if command == "build" {
		goExecutable := flags.String("go", "", "existing Go executable path; defaults to PATH")
		source := flags.String("source", ".", "evaluated checkout")
		directory := flags.String("bin-dir", "", "new input directory whose parent exists")
		budgetPath := flags.String("budget", "", "optional complete budget JSON")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *directory == "" || flags.NArg() != 0 {
			return errors.New("build requires --bin-dir and accepts no positional input")
		}
		budget := clievaluate.DefaultBudget()
		if *budgetPath != "" {
			budget = clievaluate.Budget{}
			if err := clievaluate.ReadInput(*budgetPath, &budget); err != nil {
				return err
			}
		}
		if err := clievaluate.BuildBinaries(ctx, *source, *directory, *goExecutable, budget); err != nil {
			return err
		}
		_, err := fmt.Fprintln(out, "Prepared Windows and Linux CLI, daemon and coordinator binaries in", *directory)
		return err
	}
	root := flags.String("out", "", "run evidence directory")
	if command == "capture" || command == "resume" {
		budgetPath := flags.String("budget", "", "optional complete budget JSON; defaults are recorded in the report")
		source := flags.String("source", ".", "evaluated source checkout")
		binaries := flags.String("bin-dir", "", "existing windows and linux executable directories")
		distro := flags.String("wsl-distro", "", "existing WSL2 distribution")
		user := flags.String("wsl-user", "", "existing unelevated Linux user")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *root == "" || *binaries == "" || *distro == "" || *user == "" {
			return errors.New("--out, --bin-dir, --wsl-distro and --wsl-user are required")
		}
		budget := clievaluate.DefaultBudget()
		if *budgetPath != "" {
			budget = clievaluate.Budget{}
			if err := clievaluate.ReadInput(*budgetPath, &budget); err != nil {
				return err
			}
		}
		if flags.NArg() != 0 {
			return errors.New("unexpected positional input to " + command)
		}
		report, err := clievaluate.Capture(ctx, clievaluate.CaptureOptions{Root: *root, Source: *source, BinDirectory: *binaries, WSLDistro: *distro, WSLUser: *user, Resume: command == "resume", Budget: budget})
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "Captured %d cases; %s. Review corpus.json, then finalize with the completed review submission.\n", len(report.Cases), report.Phase)
		return err
	}
	review := flags.String("reviews", "", "agent review JSON against the frozen corpus (finalize only)")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if *root == "" {
		return errors.New("--out is required")
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional input to " + command)
	}
	if command == "finalize" && *review == "" {
		return errors.New("finalize requires --reviews")
	}
	if command != "finalize" && command != "validate" && command != "render" {
		return errors.New("unknown coordinator command: " + command)
	}
	if command != "finalize" && *review != "" {
		return errors.New("--reviews is accepted only by finalize")
	}
	var report clievaluate.Report
	var err error
	if command == "validate" {
		report, err = clievaluate.ValidateRun(*root)
	} else {
		report, err = clievaluate.ProcessReport(*root, *review)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s: %s (%s), %d planned cases\n", command, report.Summary.Outcome, report.Phase, report.Summary.Planned)
	return err
}
