package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var logsCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "logs [service]",
	Short:         "short logs command",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		tail, err := flags.GetString("tail")
		if err != nil {
			return err
		}

		skipStdout, err := flags.GetBool("no-stdout")
		if err != nil {
			return err
		}

		skipStderr, err := flags.GetBool("no-stderr")
		if err != nil {
			return err
		}

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
			tasks.ContainerLogs(api, tail, skipStdout, skipStderr),
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	flags := logsCmd.Flags()
	flags.StringP("tail", "t", "all", "number of lines to show from the end of the logs for a service")
	flags.BoolP("no-stderr", "E", false, "don't show stderr output")
	flags.BoolP("no-stdout", "O", false, "don't show stdout output")
}
