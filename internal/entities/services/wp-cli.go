package services

import (
	"wpld/internal/entities"
	"wpld/internal/stdout"
)

func NewWpCliService(project entities.Project, cmd []string) entities.Service {
	wpcli := entities.Service{
		ID:           project.GetContainerIDForService("wp-cli"),
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
