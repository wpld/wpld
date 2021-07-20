package tasks

import (
	"context"
	"errors"
	"os"

	"github.com/mattn/go-isatty"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
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

		tty := false
		stdout := os.Stdout.Fd()
		if isatty.IsTerminal(stdout) {
			tty = true
		} else if isatty.IsCygwinTerminal(stdout) {
			tty = true
		}

		wpcli := entities.Service{
			ID:           project.GetContainerIDForService("wp-cli"),
			Network:      project.GetNetworkName(),
			Project:      project.Name,
			AttachStdout: true,
			AttachStdin:  true,
			AttachStderr: true,
			Tty:          tty,
			Spec: entities.Specification{
				Image: "wordpress:cli",
				Cmd:   append([]string{"wp"}, args...),
				VolumesFrom: []string{
					project.GetContainerIDForService("wp"),
				},
				Env: wp.Env,
			},
		}

		if err := api.RunContainer(ctx, wpcli); err != nil {
			return err
		}

		return next(ctx)
	}
}
