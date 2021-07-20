package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var newCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "new",
	Short:         "new short desc",
	RunE: func(c *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.NewProjectPromptPipe(),
			tasks.ProjectStructurePipe(fs),
			tasks.ProjectNginxConfigPipe(fs),
			tasks.NewProjectWPCLIConfigPipe(fs),
			tasks.ProjectMarshalPipe(fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
