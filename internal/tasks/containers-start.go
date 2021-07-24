package tasks

import (
	"context"
	"fmt"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func StartContainersPipe(api docker.Docker, pull bool) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		if err := api.EnsureNetworkExists(ctx, project.GetNetworkName()); err != nil {
			return err
		}

		for _, volume := range project.Volumes {
			if err := api.EnsureVolumeExists(ctx, volume); err != nil {
				return err
			}
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		for _, service := range services {
			if service.Spec.Name != "" {
				msg := fmt.Sprintf("%s started", service.Spec.Name)
				stdout.StartSpinner(msg)
			}

			err := api.StartContainer(ctx, service, pull)
			stdout.StopSpinner()

			if err != nil {
				return err
			}

			if service.Spec.Name != "" {
				msg := fmt.Sprintf("%s started", service.Spec.Name)
				stdout.Success(msg)
			}
		}

		return next(ctx)
	}
}
