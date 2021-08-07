package tasks

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

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

		var project entities.Project
		if err := config.Unmarshal(&project); err != nil {
			return err
		}

		if err := os.Chdir(filepath.Dir(config.ConfigFileUsed())); err != nil {
			return err
		}

		return next(context.WithValue(ctx, "project", project))
	}
}
