package tasks

import (
	"context"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ProjectUnmarshalPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		data, err := afero.ReadFile(fs, ".wpld.yml")
		if err != nil {
			return err
		}

		var project entities.Project
		if err = yaml.Unmarshal(data, &project); err != nil {
			return err
		}

		return next(context.WithValue(ctx, "project", project))
	}
}
