package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configListCmd = &cobra.Command{
	SilenceUsage: true,
	Args: cobra.NoArgs,
	Use: "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, key := range viper.AllKeys() {
			logrus.Infof("%s: %s", key, viper.GetString(key))
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
}
