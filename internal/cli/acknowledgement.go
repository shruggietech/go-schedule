package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// printAcknowledgement preserves human text and the selected structured output contract.
func printAcknowledgement(cmd *cobra.Command, value any, format string, args ...any) error {
	if jsonOutput(cmd) {
		return printJSONTo(cmd.OutOrStdout(), value)
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), format, args...)
	return err
}
