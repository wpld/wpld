package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var downCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "down",
	Short:         "down short desc",
	Aliases: []string{
		"stop",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		var pipeline pipelines.Pipeline

		if all {
			pipeline = pipelines.NewPipeline(
				tasks.StopAllContainersPipe(api),
			)
		} else {
			pipeline = pipelines.NewPipeline(
				tasks.ProjectUnmarshalPipe(fs),
				tasks.StopContainersPipe(api),
				tasks.ReloadProxyPipe(api, fs),
			)
		}

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().BoolP("all", "a", false, "stop all running projects")
}
