package specs

import (
	"wpld/internal/entities"
)

func NewNginxSpec() entities.Specification {
	return entities.Specification{
		Name:    "Nginx",
		Image:   "nginx:alpine",
		Volumes: []string{},
		VolumesFrom: []string{
			"wp",
		},
		DependsOn: []string{
			"wp",
		},
	}
}
