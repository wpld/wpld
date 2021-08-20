package tasks

import (
	"context"
	"errors"
	"os"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func WPCLIRunPipe(api docker.Docker, args []string) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		wp, ok := project.Services["wp"]
		if !ok {
			return errors.New("wp service is not defined")
		}

		wpcli := entities.Service{
			ID:           project.GetContainerIDForService("wp-cli"),
			Project:      project.Name,
			AttachStdout: true,
			AttachStdin:  true,
			AttachStderr: true,
			Tty:          stdout.IsTerm(),
			Spec: entities.Specification{
				Image: "wordpress:cli",
				Cmd:   append([]string{"wp"}, args...),
				VolumesFrom: []string{
					project.GetContainerIDForService("wp"),
				},
				Env: wp.Env,
			},
			Network: project.GetNetwork(),
		}

		if err := api.ContainerStart(ctx, wpcli, false); err != nil {
			return err
		}

		statusCode, statusErr := api.ContainerAttach(ctx, wpcli)
		if statusErr == nil {
			os.Exit(statusCode)
		}

		return statusErr
	}
}
