package tasks

import (
	"context"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ServiceFindPipe(serviceID string) pipelines.Pipe {
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

		for _, service := range services {
			if service.ID == id {
				ctx = context.WithValue(ctx, "service", service)
				break
			}
		}

		return next(ctx)
	}
}
