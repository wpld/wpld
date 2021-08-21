package tasks

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

//go:embed embeds/wordpress/wp-cli.yml
var wpcliConfig string

func ProjectWPCLIConfigPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		wp, ok := project.Services["wp"]
		if !ok {
			return next(ctx)
		}

		configFilename := ".wpld/wordpress/wp-cli.yml"
		if err := fs.MkdirAll(filepath.Dir(configFilename), 0755); err != nil {
			return err
		}

		if err := afero.WriteFile(fs, configFilename, []byte(wpcliConfig), 0644); err != nil {
			return err
		}

		wp.Volumes = append(
			wp.Volumes,
			fmt.Sprintf("%s:/var/www/html/wp-cli.yml:ro", configFilename),
		)

		project.Services["wp"] = wp

		return next(context.WithValue(ctx, "project", project))
	}
}
