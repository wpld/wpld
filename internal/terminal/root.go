package terminal

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"wpld/internal/misc"
	"wpld/internal/stdout"
)

var rootCmd = &cobra.Command{
	Use:     "wpld",
	Short:   "WordPress Local Docker",
	Version: misc.VERSION,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		pflags := cmd.Root().PersistentFlags()
		if verbosity, err := pflags.GetCount("verbose"); err == nil {
			stdout.SetLogLevel(verbosity)
		}
	},
}

func init() {
	pflags := rootCmd.PersistentFlags()
	pflags.CountP("verbose", "v", "verbose output")
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		stdout.Error(err.Error())
		os.Exit(1)
	}
}
