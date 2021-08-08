package tasks

import (
	"context"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ContainerExecPipe(api docker.Docker, cmd []string, wd string) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		service, ok := ctx.Value("service").(entities.Service)
		if !ok {
			return ServiceNotFoundErr
		}

		// TODO: implement EXEC command as it is done in the docker-cli https://github.com/docker/cli/blob/master/cli/command/container/exec.go
		if err := api.ContainerStart(ctx, service, false); err != nil {
			return err
		}

		if err := api.ContainerExecAttach(ctx, service, cmd, wd); err != nil {
			return err
		}

		return next(ctx)
	}
}
