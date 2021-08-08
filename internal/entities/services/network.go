package services

import (
	"wpld/internal/entities"
)

var globalNetwork = entities.Network{
	Name:   "wpld",
	Subnet: "10.0.1.0/24",
}
