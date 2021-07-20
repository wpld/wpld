package tasks

import (
	"context"
	"errors"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ContainerLogs(api docker.Docker, serviceID, tail string, skipStdout, skipStderr bool) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		id := project.GetContainerIDForService(serviceID)
		found := false

		for _, service := range services {
			if service.ID != id {
				continue
			} else {
				found = true
			}

			if err := api.ContainerLogs(ctx, service, tail, skipStdout, skipStderr); err != nil {
				return err
			}

			break
		}

		if !found {
			return errors.New("service not found")
		}

		return next(ctx)
	}
}
