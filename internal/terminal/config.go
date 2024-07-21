package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var configCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "config",
	Short:         "Displays configuration file",
	RunE: func(c *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.ProjectDisplayPipe(c.OutOrStdout()),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
