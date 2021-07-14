package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/cases"
	"wpld/internal/connectors/docker"
	"wpld/internal/pipelines"
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
			cases.ProjectUnmarshalPipe(fs),
			cases.StartContainersPipe(api, false), // TODO: replace "false" with the "--pull" flag value
			cases.ReloadProxyPipe(api, fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
