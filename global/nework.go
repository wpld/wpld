package global

import (
	"wpld/models"
)

const (
	NETWORK_NAME = "wpld"
)

func VerifyNetwork(factory models.DockerFactory) error {
	return factory.Network(NETWORK_NAME).Create()
}
