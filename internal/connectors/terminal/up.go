package terminal

import (
	"github.com/spf13/cobra"

	"wpld/internal/cases"
)

var upCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "up",
	Short:         "up short desc",
	Aliases: []string{
		"start",
	},
	RunE: func(c *cobra.Command, args []string) error {
		return cases.StartProjectPipeline(fs).Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
