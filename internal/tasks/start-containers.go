package tasks

import (
	"context"
	"errors"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func StartContainersPipe(api docker.Docker, pull bool) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return errors.New("project not found")
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
			if err := api.StartContainer(ctx, service, pull); err != nil {
				return err
			}
		}

		return next(ctx)
	}
}
