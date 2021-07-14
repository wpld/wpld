package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/connectors/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var downCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "down",
	Aliases: []string{
		"stop",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.StopContainersPipe(api),
			tasks.ReloadProxyPipe(api, fs),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
