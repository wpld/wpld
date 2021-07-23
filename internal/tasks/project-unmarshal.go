package tasks

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ProjectUnmarshalPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		var config string

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		for {
			config = filepath.Join(dir, ".wpld.yml")
			if exists, err := afero.Exists(fs, config); err != nil {
				return err
			} else if exists {
				break
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			} else {
				dir = parent
			}
		}

		data, err := afero.ReadFile(fs, config)
		if err != nil {
			return err
		}

		var project entities.Project
		if err = yaml.Unmarshal(data, &project); err != nil {
			return err
		}

		if err := os.Chdir(dir); err != nil {
			return err
		}

		return next(context.WithValue(ctx, "project", project))
	}
}
