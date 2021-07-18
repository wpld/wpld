package terminal

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"wpld/internal/misc"
)

var rootCmd = &cobra.Command{
	Use:     "wpld",
	Short:   "wpld short desc",
	Version: misc.VERSION,
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
}
