package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/shruggietech/go-schedule/internal/bundle"
)

func newBundleCmd() *cobra.Command {
	command := &cobra.Command{Use: "bundle", Short: "Export, review, and apply portable automation intent"}
	command.AddCommand(bundleExport(), bundleValidate(), bundleCompare(), bundlePreview(), bundleApply())
	return command
}

func bundleExport() *cobra.Command {
	var output string
	command := &cobra.Command{Use: "export", Short: "Export deterministic, secret-free automation intent", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := reqCtx()
		defer cancel()
		document, err := newClient().ExportBundle(ctx)
		if err != nil {
			return err
		}
		return writeBundleJSON(cmd, output, document)
	}}
	command.Flags().StringVarP(&output, "output", "o", "", "write JSON to a file")
	return command
}

func bundleValidate() *cobra.Command {
	return &cobra.Command{Use: "validate <bundle-file>", Short: "Validate a portable bundle without mutation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		document, err := readBundleDocument(args[0])
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx()
		defer cancel()
		result, err := newClient().ValidateBundle(ctx, document)
		if err != nil {
			return err
		}
		if err := printJSONTo(cmd.OutOrStdout(), result); err != nil {
			return err
		}
		if !result.Valid {
			return fmtUsage("bundle validation failed")
		}
		return nil
	}}
}

func bundlePreview() *cobra.Command {
	var output string
	command := &cobra.Command{Use: "preview <bundle-file>", Short: "Create a target-bound, single-use apply plan", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		document, err := readBundleDocument(args[0])
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx()
		defer cancel()
		plan, err := newClient().PreviewBundle(ctx, document)
		if err != nil {
			return err
		}
		return writeBundleJSON(cmd, output, plan)
	}}
	command.Flags().StringVarP(&output, "output", "o", "", "write the single-use plan to a file")
	return command
}

func bundleCompare() *cobra.Command {
	return &cobra.Command{Use: "compare <bundle-file>", Short: "Compare a bundle with the selected target without mutation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		document, err := readBundleDocument(args[0])
		if err != nil {
			return err
		}
		ctx, cancel := reqCtx()
		defer cancel()
		plan, err := newClient().CompareBundle(ctx, document)
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), plan)
	}}
}

func bundleApply() *cobra.Command {
	return &cobra.Command{Use: "apply <plan-file>", Short: "Apply a reviewed, target-bound plan", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		var plan bundle.Plan
		if err := readBundleJSON(args[0], &plan); err != nil {
			return err
		}
		ctx, cancel := reqCtx()
		defer cancel()
		result, err := newClient().ApplyBundle(ctx, plan)
		if err != nil {
			return err
		}
		return printJSONTo(cmd.OutOrStdout(), result)
	}}
}

func readBundleDocument(path string) (bundle.Document, error) {
	var document bundle.Document
	if err := readBundleJSON(path, &document); err != nil {
		return bundle.Document{}, err
	}
	return document, nil
}

func readBundleJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("read bundle JSON: %w", err)
	}
	return nil
}

func writeBundleJSON(cmd *cobra.Command, path string, value any) error {
	if path == "" {
		return printJSONTo(cmd.OutOrStdout(), value)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}
