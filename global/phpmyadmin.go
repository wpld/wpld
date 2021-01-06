package global

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"wpld/utils"
)

const (
	MYADMIN_IMAGE_NAME = "phpmyadmin:latest"
	MYADMIN_CONTAINER_NAME = "wpld_global_phpmyadmin"
)

func RunMyAdmin(ctx context.Context, cli *client.Client, pull bool) error {
	img := utils.Image{
		Name: MYADMIN_IMAGE_NAME,
	}

	if pull {
		if err := img.Pull(ctx, cli); err != nil {
			return err
		}
	}

	myadmin := utils.Container{
		Name: MYADMIN_CONTAINER_NAME,
		Create: &container.Config{
			Hostname: "phpmyadmin",
			Image: img.Name,
			Env: []string{
				"PMA_HOST=" + MYSQL_CONTAINER_NAME,
				"UPLOAD_LIMIT=512MiB",
			},
		},
		Host: &container.HostConfig{
			NetworkMode: NETWORK_NAME,
			IpcMode: "shareable",
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP: "127.0.0.1",
						HostPort: "8092",
					},
				},
			},
		},
	}

	return myadmin.Start(ctx, cli)
}

func StopMyAdmin(ctx context.Context, cli *client.Client, rm bool) error {
	myadmin := utils.Container{
		Name: MYADMIN_CONTAINER_NAME,
	}

	if rm {
		return myadmin.Remove(ctx, cli)
	}

	return myadmin.Stop(ctx, cli)
}
