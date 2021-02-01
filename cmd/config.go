package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Use:          "config",
	Short:        "A brief description of your command",
}

func init() {
	configCmd.Long = heredoc.Doc(`
		A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.
	`)

	rootCmd.AddCommand(configCmd)
}
