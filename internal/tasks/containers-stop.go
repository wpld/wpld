package tasks

import (
	"context"
	"fmt"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func StopContainersPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		for _, service := range services {
			if service.Spec.Name != "" {
				msg := fmt.Sprintf("Stopping %s...", service.Spec.Name)
				stdout.StartSpinner(msg)
			}

			err := api.StopContainer(ctx, service)
			stdout.StopSpinner()

			if err != nil {
				return err
			}

			if service.Spec.Name != "" {
				msg := fmt.Sprintf("%s stopped", service.Spec.Name)
				stdout.Success(msg)
			}
		}

		return next(ctx)
	}
}

func StopAllContainersPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		if err := api.StopAllContainers(ctx); err != nil {
			return err
		}

		return next(ctx)
	}
}
