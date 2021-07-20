package tasks

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ProjectStructurePipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		if err := fs.MkdirAll(filepath.Join(project.ID, ".wpld"), 0755); err != nil {
			return err
		}

		if err := os.Chdir(project.ID); err != nil {
			return err
		}

		return next(ctx)
	}
}
