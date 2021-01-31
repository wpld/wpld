package cases

import (
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/viper"
	"strings"
	"wpld/global"
	"wpld/models"
	"wpld/utils"
)

func StartArbitraryContainer(factory models.DockerFactory, service *viper.Viper, pull bool) error {
	config := &container.Config{
		Image: service.GetString("image"),
	}

	img := factory.Image(config.Image)
	// TODO: pull image if it doesn't exist
	if pull {
		if err := img.Pull(); err != nil {
			return err
		}
	}

	host := &container.HostConfig{
		NetworkMode: global.NETWORK_NAME,
		IpcMode:     "shareable",
	}

	volumes := service.GetStringSlice("volumes")
	if len(volumes) > 0 {
		host.Binds = []string{}
		for _, volume := range volumes {
			host.Binds = append(host.Binds, utils.NormalizeVolumeBind(volume))
		}
	}

	env := service.GetStringMapString("env")
	if len(env) > 0 {
		config.Env = []string{}
		for key, value := range env {
			config.Env = append(config.Env, strings.ToUpper(key)+"="+value)
		}
	}

	cntr := factory.Container(service.GetString("name"))
	if err := cntr.Create(config, host); err != nil {
		return err
	}

	return cntr.Start()
}
