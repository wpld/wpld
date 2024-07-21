package tasks

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ProjectDisplayPipe(w io.Writer) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		buffer := bytes.NewBufferString("")

		encoder := yaml.NewEncoder(buffer)
		encoder.SetIndent(2)

		if err := encoder.Encode(project); err != nil {
			return err
		}

		fmt.Fprint(w, buffer.String())

		return next(ctx)
	}
}
