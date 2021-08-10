package tasks

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ProjectUnmarshalPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		config := viper.New()
		config.SetFs(fs)
		config.SetConfigName(".wpld")
		config.SetConfigType("yaml")
		config.AddConfigPath(".")

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		for {
			if parent := filepath.Dir(dir); parent != dir {
				config.AddConfigPath(parent)
				dir = parent
			} else {
				break
			}
		}

		if err := config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return ProjectNotFoundErr
			} else {
				return err
			}
		}

		// Viper updates all env variables to be in a lower case register. This is not acceptable for us, so we need to
		// read the config file again and parse it manually.
		configFile := config.ConfigFileUsed()
		data, err := afero.ReadFile(fs, configFile)
		if err != nil {
			return err
		}

		var project entities.Project
		if err := yaml.Unmarshal(data, &project); err != nil {
			return err
		}

		if err := os.Chdir(filepath.Dir(configFile)); err != nil {
			return err
		}

		return next(context.WithValue(ctx, "project", project))
	}
}
