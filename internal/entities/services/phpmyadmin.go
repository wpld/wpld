package services

import (
	"wpld/internal/entities"
)

func NewPHPMyAdminService() entities.Service {
	return entities.Service{
		ID:      "wpld__phpmyadmin",
		Network: GetGlobalNetwork(),
		Spec: entities.Specification{
			Image: "phpmyadmin:latest",
			Env: map[string]string{
				"UPLOAD_LIMIT": "1024MiB",
			},
			Domains: []string{
				"phpmyadmin",
				"phpmyadmin.test",
			},
		},
	}
}
