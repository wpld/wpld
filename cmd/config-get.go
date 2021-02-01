package cmd

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configGetCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	Use:          "get [key]",
	Short:        "A brief description of your command",
	RunE:         runConfigGet,
}

func init() {
	configGetCmd.Long = heredoc.Doc(`
		A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.
	`)

	configCmd.AddCommand(configGetCmd)
}

func runConfigGet(_ *cobra.Command, args []string) error {
	fmt.Println(viper.GetString(args[0]))
	return nil
}
