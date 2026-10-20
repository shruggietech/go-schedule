package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRuntimeCmd() *cobra.Command {
	return &cobra.Command{Use: "runtime", Short: "Show the selected local daemon's effective storage and process metadata", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		api := clientFor(cmd)
		if api.Remote() {
			return fmtUsage("runtime metadata is available only through local IPC; select the local scheduler or use --ipc-path")
		}
		ctx, cancel := reqCtx(cmd)
		defer cancel()
		info, err := api.RuntimeInfo(ctx)
		if err != nil {
			return err
		}
		if jsonOutput(cmd) {
			return printJSONTo(cmd.OutOrStdout(), info)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Data directory: %s\nDatabase: %s\nConfiguration: %s\nLog: %s\nLock: %s\nIPC endpoint: %s\nProcess ID: %d\n", info.DataDir, info.DatabasePath, displayOrDash(info.ConfigPath), info.LogPath, info.LockPath, displayOrDash(info.IPCPath), info.ProcessID)
		return err
	}}
}
