package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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
		config := viper.AllSettings()

		if data, err := yaml.Marshal(config); err != nil {
			logrus.Errorf("Unable to marshal config to YAML: %v", err)
		} else {
			fmt.Print(string(data))
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
}
