package tasks

import (
	"context"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

func NetworksCreatePipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		if err := api.EnsureNetworkExists(ctx, project.GetNetwork()); err != nil {
			return err
		}

		if err := api.EnsureNetworkExists(ctx, services.GetGlobalNetwork()); err != nil {
			return err
		}

		return next(ctx)
	}
}
