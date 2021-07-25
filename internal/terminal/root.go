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
	Short:   "wpld short desc",
	Version: misc.VERSION,
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		stdout.Error(err.Error())
		os.Exit(1)
	}
}
