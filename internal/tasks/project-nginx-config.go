package tasks

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

//go:embed embeds/nginx/default.conf
var nginxConfig string

func ProjectNginxConfigPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return errors.New("project is not found")
		}

		if nginx, ok := project.Services["nginx"]; ok {
			configFilename := ".wpld/nginx/default.conf.template"

			if err := fs.MkdirAll(filepath.Dir(configFilename), 0755); err != nil {
				return err
			}

			if err := afero.WriteFile(fs, configFilename, []byte(nginxConfig), 0644); err != nil {
				return err
			}

			if nginx.Env == nil {
				nginx.Env = make(map[string]string)
			}

			nginx.Env["PHPFPM_HOST"] = "wp"
			nginx.Env["PHPFPM_PORT"] = "9000"

			nginx.Volumes = append(
				nginx.Volumes,
				fmt.Sprintf("%s:/etc/nginx/templates/default.conf.template:ro", configFilename),
			)

			project.Services["nginx"] = nginx

			return next(context.WithValue(ctx, "project", project))
		}

		return next(ctx)
	}
}
