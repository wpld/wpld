package services

import (
	"fmt"

	"wpld/internal/entities"
	"wpld/internal/stdout"
)

var (
	wpcliIndex int = 0
)

func NewWpCliService(project entities.Project, cmd []string) entities.Service {
	wpcliIndex++

	id := fmt.Sprintf("wp-cli_%d", wpcliIndex)
	wpcli := entities.Service{
		ID:           project.GetContainerIDForService(id),
		Project:      project.Name,
		AttachStdout: true,
		AttachStdin:  true,
		AttachStderr: true,
		Tty:          stdout.IsTerm(),
		Network:      project.GetNetwork(),
		Spec: entities.Specification{
			Image: "wordpress:cli",
			Cmd:   cmd,
			VolumesFrom: []string{
				project.GetContainerIDForService("wp"),
			},
		},
	}

	if wp, ok := project.Services["wp"]; ok {
		wpcli.Spec.Env = wp.Env
	}

	return wpcli
}
