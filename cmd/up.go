package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	SilenceUsage: true,
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("up called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
