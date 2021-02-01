package cmd

import (
	"errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"wpld/config"
	"wpld/utils"
)

var (
	configFilename string
	rootCmd        = &cobra.Command{
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		Use:           "wpld",
		Short:         "short desc",
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
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	rootCmd.Long = heredoc.Doc(`
		long desc
	`)

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
				Code:            utils.HOMEDIR_DETECTION_ERROR,
				FriendlyMessage: "Can't find the current user home directory.",
				OriginalError:   err,
			})
		} else {
			viper.AddConfigPath(home)
			viper.SetConfigName(".wpld")
			viper.SetConfigType("yaml")
		}
	}

	viper.AutomaticEnv()
	viper.SetFs(afero.NewOsFs())

	viper.SetDefault("log.level", "info")

	viper.SetDefault(config.MYSQL_PORT, "3306")
	viper.SetDefault(config.MYSQL_MEMORY, "256MiB")
	viper.SetDefault(config.MYSQL_RESERVATION, "256MiB")

	viper.SetDefault(config.PHPMYADMIN_PORT, "8092")

	_ = viper.SafeWriteConfig()
	if err := viper.ReadInConfig(); err != nil {
		must(utils.ExecutionError{
			Code:            utils.CONFIG_ERROR,
			FriendlyMessage: "Can't read the config file.",
			OriginalError:   err,
		})
	}
}

func initLogger() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetLevel(level)
	}
}
