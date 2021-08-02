package services

import (
	"wpld/internal/entities"
)

func NewPHPMyAdminService() entities.Service {
	return entities.Service{
		ID:      "wpld__phpmyadmin",
		Network: "bridge",
		Spec: entities.Specification{
			Image: "phpmyadmin:fpm-alpine",
			Ports: []string{
				"80:9000",
			},
			Env: map[string]string{
				"UPLOAD_LIMIT": "1024MiB",
			},
			Domains: []string{
				"phpmyadmin.local",
			},
		},
	}
}
