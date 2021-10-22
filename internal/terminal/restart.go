package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var restartCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "restart",
	Short:         "restart short desc",
	RunE: func(cmd *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.ContainersStopPipe(api),
			tasks.ContainersStartPipe(api, false),
			//tasks.DNSReloadPipe(api, fs),
			tasks.GlobalProxyReload(api, fs),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
