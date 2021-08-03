package tasks

import (
	"context"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ContainerSSHPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		service, ok := ctx.Value("service").(entities.Service)
		if !ok {
			return ServiceNotFoundErr
		}

		if err := api.ContainerStart(ctx, service, false); err != nil {
			return err
		}

		if err := api.ContainerAttach(ctx, service); err != nil {
			return err
		}

		return next(ctx)
	}
}