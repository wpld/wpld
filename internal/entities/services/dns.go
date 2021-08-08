package services

import (
	"wpld/internal/entities"
)

func NewDnsService() entities.Service {
	return entities.Service{
		ID:      "wpld__dnsmasq",
		Network: globalNetwork,
		Spec: entities.Specification{
			Image:     "4km3/dnsmasq:latest",
			IPAddress: "10.0.1.53",
			CapAdd: []string{
				"NET_ADMIN",
			},
		},
	}
}
