package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var startCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "start",
	Short:         "start short desc",
	Aliases: []string{
		"up",
	},
	RunE: func(c *cobra.Command, args []string) error {
		pull, err := c.Flags().GetBool("pull")
		if err != nil {
			return err
		}

		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.NetworksCreatePipe(api),
			tasks.ContainersStartPipe(api, pull),
			tasks.PHPMyAdminReloadPipe(api),
			//tasks.DNSReloadPipe(api, fs),
			tasks.GlobalProxyReload(api, fs),
			// tasks.WordPressInstallPipe(api),
			tasks.ProjectInformationPipe(api),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("pull", "p", false, "force pulling images before starting containers")
}
