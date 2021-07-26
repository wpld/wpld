package tasks

import (
	"context"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ContainerLogsPipe(api docker.Docker, tail string, skipStdout, skipStderr bool) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		service, ok := ctx.Value("service").(entities.Service)
		if !ok {
			return ServiceNotFoundErr
		}

		if err := api.ContainerLogs(ctx, service, tail, skipStdout, skipStderr); err != nil {
			return err
		}

		return next(ctx)
	}
}
