package tasks

import (
	"context"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

func NetworksRemovePipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		network := project.GetNetworkName()
		if isUsed, isUsedErr := api.NetworkIsInUsed(ctx, network); !isUsed && isUsedErr == nil {
			if err := api.NetworkRemove(ctx, network); err != nil {
				return err
			}
		}

		network = services.GetGlobalNetwork().Name
		if isUsed, isUsedErr := api.NetworkIsInUsed(ctx, network); !isUsed && isUsedErr == nil {
			if err := api.NetworkRemove(ctx, network); err != nil {
				return err
			}
		}

		return next(ctx)
	}
}
