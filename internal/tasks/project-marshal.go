package tasks

import (
	"bytes"
	"context"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"wpld/internal/pipelines"
)

func ProjectMarshalPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		buffer := bytes.NewBufferString("")

		encoder := yaml.NewEncoder(buffer)
		encoder.SetIndent(2)

		if err := encoder.Encode(ctx.Value("project")); err != nil {
			return err
		}

		if err := afero.WriteFile(fs, ".wpld.yml", buffer.Bytes(), 0644); err != nil {
			return err
		}

		return next(ctx)
	}
}
