package cmd

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"wpld/utils"
)

var (
	configFilename string
	rootCmd = &cobra.Command{
		SilenceErrors: true,
		Use: "wpld",
		Short: "short desc",
		Long: "long desc",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.Debugf("Running {%s} command...", cmd.Use)
		},
	}
)

func must(err error) {
	if err != nil {
		exitCode := utils.UKNOWN_ERROR

		var execErr utils.ExecutionError
		if errors.As(err, &execErr) && execErr.Code > 0 {
			exitCode = execErr.Code
		}

		logrus.Error(err)
		logrus.Exit(exitCode)
	}
}

func Execute() {
	must(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// TODO: level should be read from the config file
	logrus.SetLevel(logrus.DebugLevel)

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
		logrus.Debugf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		// TODO: create a new default config if it doesn't exist yet
	}
}
