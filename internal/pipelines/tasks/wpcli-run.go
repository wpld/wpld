package tasks

import (
	"context"
	"os"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

func WPCLIRunPipe(api docker.Docker, args []string) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		if _, wp := project.Services["wp"]; !wp {
			return WpServiceNotFoundErr
		}

		wpcli := services.NewWpCliService(project, append([]string{"wp"}, args...))
		if code, err := api.ContainerRun(ctx, wpcli); err == nil {
			os.Exit(code)
		} else {
			return err
		}

		return next(ctx)
	}
}
