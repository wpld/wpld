package terminal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wpld",
	Short: "wpld short desc",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
