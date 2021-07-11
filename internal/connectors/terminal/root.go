package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var fs = afero.NewOsFs()

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
