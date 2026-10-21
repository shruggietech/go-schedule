package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
)

func writeHumanTable(output io.Writer, headings string, rows func(io.Writer) error) error {
	table := tabwriter.NewWriter(output, 0, 2, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, headings); err != nil {
		return err
	}
	if err := rows(table); err != nil {
		return err
	}
	return table.Flush()
}
