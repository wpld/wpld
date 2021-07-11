package terminal

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wpld",
	Short: "wpld short desc",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
