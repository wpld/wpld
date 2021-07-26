package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var newCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "new",
	Short:         "new short desc",
	RunE: func(c *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectPromptPipe(),
			tasks.ProjectStructurePipe(fs),
			tasks.ProjectNginxConfigPipe(fs),
			tasks.ProjectPHPConfigPipe(fs),
			tasks.ProjectWPCLIConfigPipe(fs),
			tasks.ProjectMarshalPipe(fs),
			tasks.ContainersStartPipe(api, false),
			tasks.ReloadProxyPipe(api, fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
