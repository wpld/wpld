package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var wpCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "wp COMMAND [ARG...]",
	Short:         "Executes a WP CLI command in the WordPress container",
	RunE: func(cmd *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.WPCLIRunPipe(api, args),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(wpCmd)
}
