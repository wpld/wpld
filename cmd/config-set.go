package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configSetCmd = &cobra.Command{
	SilenceUsage: true,
	Args: cobra.ExactArgs(2),
	Use: "set [key] [value]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set(args[0], args[1])
		return viper.WriteConfig()
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
