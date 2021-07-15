package services

import (
	"wpld/internal/entities"
)

func NewProxyService() entities.Service {
	return entities.Service{
		ID:      "wpld__reverse_proxy",
		Network: "host",
		Spec: entities.Specification{
			Image: "nginx:alpine",
		},
	}
}
