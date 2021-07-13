package cases

import (
	"context"
	"errors"

	"wpld/internal/connectors/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func StopContainersPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return errors.New("project not found")
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		for _, service := range services {
			if err := api.StopContainer(ctx, service); err != nil {
				return err
			}
		}

		return next(ctx)
	}
}
