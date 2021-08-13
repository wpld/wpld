package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/pipelines/tasks"
)

var execCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "exec [SERVICE] [COMMAND] [ARG...]",
	Short:         "exec short desc",
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

		command := []string{"bash"}
		if len(args) > 1 {
			command = args[1:]
		}

		wd, err := cmd.Flags().GetString("working-dir")
		if err != nil {
			return err
		}

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.ServiceFindPipe(service),
			tasks.ContainerExecPipe(api, command, wd),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringP("working-dir", "w", "", "working directory")
}
