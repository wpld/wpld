package terminal

import (
	"github.com/spf13/cobra"

	"wpld/internal/pipelines"
)

var logsCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "logs",
	Short:         "short logs command",
	RunE: func(cmd *cobra.Command, args []string) error {
		pipeline := pipelines.NewPipeline(
		// TODO: Add pipes to implement the logs command
		)

		return pipeline.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	flags := logsCmd.Flags()
	flags.StringP("service", "s", "", "service name to get logs for")
	flags.IntP("tail", "t", 0, "number of lines to show from the end of the logs for a service")
}
