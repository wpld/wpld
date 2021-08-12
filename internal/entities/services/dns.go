package services

import (
	"wpld/internal/entities"
)

func NewDnsService() entities.Service {
	return entities.Service{
		ID:      "wpld__dnsmasq",
		Network: GetGlobalNetwork(),
		Spec: entities.Specification{
			Image:     "4km3/dnsmasq:latest",
			IPAddress: "10.1.0.53",
			CapAdd: []string{
				"NET_ADMIN",
			},
		},
	}
}
