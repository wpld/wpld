package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/cases"
	"wpld/internal/connectors/docker"
	"wpld/internal/pipelines"
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
			cases.ProjectUnmarshalPipe(fs),
			cases.StopContainersPipe(api),
			cases.ReloadProxyPipe(api),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
