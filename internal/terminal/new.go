package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var newCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "new",
	Short:         "Creates a new project",
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
			tasks.NetworksCreatePipe(api),
			tasks.ContainersStartPipe(api, false),
			tasks.PHPMyAdminReloadPipe(api),
			tasks.DNSReloadPipe(api, fs),
			// tasks.WordPressInstallPipe(api),
			tasks.ProjectInformationPipe(api),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
