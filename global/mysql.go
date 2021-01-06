package global

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

	mysql := utils.Container{
		Name: MYSQL_CONTAINER_NAME,
		Create: &container.Config{
			Image: img.Name,
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
