package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var sshCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "ssh [service]",
	Short:         "ssh short desc",
	RunE: func(cmd *cobra.Command, args []string) error {
		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		service := "wp"
		if len(args) > 0 {
			service = args[0]
		}

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.ServiceFindPipe(service),
			tasks.ContainerSSHPipe(api),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
