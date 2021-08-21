package tasks

import (
	"context"
	"fmt"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func ContainersStopPipe(api docker.Docker) pipelines.Pipe {
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

			err := api.ContainerStop(ctx, service)
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

func ContainersStopAllPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		stdout.Infoln("Looking for docker containers...")
		list, err := api.FindAllRunningContainers(ctx)

		if err != nil {
			return err
		} else {
			stdout.Infof("Found %d containers", len(list))
		}

		for _, container := range list {
			stdout.Infof("Stopping %s (%s) container...", container.Names[0], container.ID[0:12])

			project, hasProjectLabel := container.Labels["io.wpld.project"]
			service, hasServiceLabel := container.Labels["io.wpld.service"]

			if hasProjectLabel && hasServiceLabel {
				msg := fmt.Sprintf("{%s} Stopping %s...", project, service)
				stdout.StartSpinner(msg)
			}

			err := api.ContainerStop(ctx, entities.Service{ID: container.ID})
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
