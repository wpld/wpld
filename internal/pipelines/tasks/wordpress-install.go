package tasks

import (
	"context"
	"os"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

func WordPressInstallPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		if _, wp := project.Services["wp"]; !wp {
			return WpServiceNotFoundErr
		}

		nginx, ok := project.Services["nginx"]
		if !ok {
			return next(ctx)
		}

		// TODO: Add progress

		cmd := []string{"wp", "core", "is-installed"}
		wpcli := services.NewWpCliService(project, cmd)
		code, err := api.ContainerRun(ctx, wpcli)
		if err != nil {
			return err
		} else if code == 0 {
			return next(ctx)
		}

		cmd = []string{
			"wp",
			"core",
			"install",
			"--url",
			// TODO: make sure domain exists
			nginx.Domains[0],
			"--title",
			project.Name,
			"--admin_user",
			project.WP.User,
			"--admin_password",
			project.WP.Password,
			"--admin_email",
			project.WP.Email,
		}

		wpcli = services.NewWpCliService(project, cmd)
		code, err = api.ContainerRun(ctx, wpcli)
		if err != nil {
			return err
		} else if code != 0 {
			os.Exit(code)
		}

		return next(ctx)
	}
}
