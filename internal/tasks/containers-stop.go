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
		list, err := api.FindAllRunningContainers(ctx)
		if err != nil {
			return err
		}

		for _, container := range list {
			project, hasProjectLabel := container.Labels["wpld.project"]
			service, hasServiceLabel := container.Labels["wpld.service"]

			if hasProjectLabel && hasServiceLabel {
				msg := fmt.Sprintf("{%s} Stopping %s...", project, service)
				stdout.StartSpinner(msg)
			}

			err := api.StopContainer(ctx, entities.Service{ID: container.ID})
			stdout.StopSpinner()

			if err != nil {
				return err
			}

			if hasProjectLabel && hasServiceLabel {
				msg := fmt.Sprintf("{%s} %s stopped", project, service)
				stdout.Success(msg)
			}
		}

		return next(ctx)
	}
}
