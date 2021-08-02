package services

import (
	"wpld/internal/entities"
)

func NewDnsService() entities.Service {
	return entities.Service{
		ID: "wpld__dnsmasq",
		Spec: entities.Specification{
			Image: "4km3/dnsmasq:latest",
			CapAdd: []string{
				"NET_ADMIN",
			},
		},
	}
}
