package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/connectors/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var upCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "up",
	Short:         "up short desc",
	Aliases: []string{
		"start",
	},
	RunE: func(c *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.StartContainersPipe(api, false), // TODO: replace "false" with the "--pull" flag value
			tasks.ReloadProxyPipe(api, fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
