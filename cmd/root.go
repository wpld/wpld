package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"wpld/utils"
)

var (
	configFilename string
	rootCmd = &cobra.Command{
		SilenceErrors: true,
		Use: "wpld",
		Short: "short desc",
		Long: "long desc",
	}
)

func must(err error) {
	if err != nil {
		exitCode := utils.UKNOWN_ERROR

		var execErr utils.ExecutionError
		if errors.As(err, &execErr) && execErr.Code > 0 {
			exitCode = execErr.Code
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode)
	}
}

func Execute() {
	must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&configFilename,
		"config",
		"",
		"Path to the configuration file. By default it uses \"$HOME/.wpld.yaml\".",
	)
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			must(utils.ExecutionError{
				Code: utils.HOMEDIR_DETECTION_ERROR,
				FriendlyMessage: "Can't find the current user home directory.",
				OriginalError: err,
			})
		} else {
			viper.AddConfigPath(home)
			viper.SetConfigName(".wpld")
		}
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// TODO: create a new config
	}
}
