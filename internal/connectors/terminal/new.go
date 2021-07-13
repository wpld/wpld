package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/cases"
	"wpld/internal/controllers/pipelines"
)

var newCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "new",
	Short:         "new short desc",
	RunE: func(c *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			cases.NewProjectPromptPipe(),
			cases.ProjectMarshalPipe(fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
