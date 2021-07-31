package services

import (
	"wpld/internal/entities"
)

func NewProxyService() entities.Service {
	return entities.Service{
		ID: "wpld__reverse_proxy",
		Spec: entities.Specification{
			Image: "4km3/dnsmasq",
			CapAdd: []string{
				"NET_ADMIN",
			},
		},
	}
}
