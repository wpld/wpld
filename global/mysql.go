package global

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	"wpld/config"
	"wpld/utils"
)

const (
	MYSQL_IMAGE_NAME = "mysql:5"
	MYSQL_CONTAINER_NAME = "wpld_global_mysql"
)

func RunMySQL(ctx context.Context, cli *client.Client, pull bool) error {
	img := utils.Image{
		Name: MYSQL_IMAGE_NAME,
	}

	if pull {
		if err := img.Pull(ctx, cli); err != nil {
			return err
		}
	}

	resources := container.Resources{}
	port := nat.PortBinding{
		HostIP: "127.0.0.1",
		HostPort: "3306",
	}

	memory := viper.GetString(config.MYSQL_MEMORY)
	if mem, err := utils.ParseBytes(memory); err == nil {
		resources.Memory = mem
	}

	reservation := viper.GetString(config.MYSQL_RESERVATION)
	if reserve, err := utils.ParseBytes(reservation); err == nil {
		resources.MemoryReservation = reserve
	}

	if viper.IsSet(config.MYSQL_PORT) {
		port.HostPort = viper.GetString(config.MYSQL_PORT)
	}

	mysql := utils.Container{
		Name: MYSQL_CONTAINER_NAME,
		Create: &container.Config{
			Image: img.Name,
			//User: strconv.Itoa(os.Getuid()),
			Env: []string{
				"MYSQL_ROOT_PASSWORD=password",
			},
		},
		Host: &container.HostConfig{
			NetworkMode: NETWORK_NAME,
			IpcMode: "shareable",
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{ port },
			},
			Resources: resources,
		},
	}

	return mysql.Start(ctx, cli)
}

func StopMySQL(ctx context.Context, cli *client.Client, rm bool) error {
	mysql := utils.Container{
		Name: MYSQL_CONTAINER_NAME,
	}

	if rm {
		return mysql.Remove(ctx, cli)
	}

	return mysql.Stop(ctx, cli)
}
