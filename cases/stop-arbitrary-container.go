package cases

import (
	"github.com/spf13/viper"
	"wpld/models"
)

func StopArbitraryContainer(factory models.DockerFactory, service *viper.Viper, rm bool) error {
	container := factory.Container(service.GetString("name"))

	if rm {
		return container.Remove()
	}

	return container.Stop()
}
