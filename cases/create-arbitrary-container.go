package cases

import (
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/viper"
	"wpld/global"
	"wpld/utils"
)

func CreateArbitraryContainer(service *viper.Viper) *utils.Container {
	create := &container.Config{
		Image: service.GetString("image"),
	}

	host := &container.HostConfig{
		NetworkMode: global.NETWORK_NAME,
		IpcMode: "shareable",
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
		create.Env = []string{}
		for key, value := range env {
			create.Env = append(create.Env, key + "=" + value)
		}
	}

	return &utils.Container{
		Name: service.GetString("name"),
		Create: create,
		Host: host,
	}
}
