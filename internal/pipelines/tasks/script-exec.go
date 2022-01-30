package tasks

import (
	"context"
	"os"
	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ScriptExecPipe(api docker.Docker, scriptID string) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		script, ok := project.Scripts[scriptID]
		if !ok {
			return ScriptNotFoundErr
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		id := project.GetContainerIDForService(script.Service)

		for _, service := range services {
			if service.ID == id {
				statusCode, statusErr := api.ContainerExecAttach(ctx, service, script.Command, script.WorkingDir)
				if statusErr == nil {
					os.Exit(statusCode)
				}

				return statusErr
			}
		}

		return next(ctx)
	}
}
