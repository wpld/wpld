package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var stopCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "stop",
	Short:         "Stop project services",
	Aliases: []string{
		"down",
		"halt",
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
				tasks.ContainersStopAllPipe(api),
				tasks.NetworksRemovePipe(api),
			)
		} else {
			pipeline = pipelines.NewPipeline(
				tasks.ProjectUnmarshalPipe(fs),
				tasks.ContainersStopPipe(api),
				tasks.PHPMyAdminReloadPipe(api),
				//tasks.DNSReloadPipe(api, fs),
				tasks.GlobalProxyReload(api, fs),
				tasks.NetworksRemovePipe(api),
			)
		}

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolP("all", "a", false, "stop all running projects")
}
