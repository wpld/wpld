package terminal

import (
	"github.com/spf13/cobra"

	"wpld/internal/cases"
)

var newCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "new",
	Short:         "new short desc",
	RunE: func(c *cobra.Command, args []string) error {
		return cases.NewProjectPipeline(fs).Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
