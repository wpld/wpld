package services

import (
	"wpld/internal/entities"
)

func GetGlobalNetwork() entities.Network {
	return entities.Network{
		Name:   "wpld",
		Subnet: "10.1.0.0/24",
	}
}
