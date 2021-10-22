package services

import (
	"wpld/internal/entities"
)

func NewProxyService() entities.Service {
	return entities.Service{
		ID:      "wpld__reverse_proxy",
		Network: GetGlobalNetwork(),
		Spec: entities.Specification{
			Image: "nginx:alpine",
			//IPAddress: "10.1.0.2",
			Ports: []string{
				"80:80",
				"443:443/tcp",
			},
		},
	}
}
