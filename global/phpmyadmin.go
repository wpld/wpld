package global

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"wpld/utils"
)

const (
	MYADMIN_IMAGE_NAME = "phpmyadmin:latest"
	MYADMIN_CONTAINER_NAME = "wpld_global_phpmyadmin"
)

func getBase64EncodedPMAConfig() string {
	servers := []map[string]string{
		{
			"host": MYSQL_CONTAINER_NAME,
			"label": "global",
		},
	}

	buffer := bytes.NewBufferString(`
<?php

$cfg['blowfish_secret'] = 'l3+wF5o$MUK@hj;[HLkQ4#V9-m?b4JmgXa]H_{uH#H]x|oQI%c1s|wFOGTc[<{3M';
$cfg['ServerDefault']   = 1;
`)

	for i, server := range servers {
		config := `
$cfg['Servers'][%[1]d]['host']      = '%[2]s';
$cfg['Servers'][%[1]d]['auth_type'] = 'config';
$cfg['Servers'][%[1]d]['user']      = 'root';
$cfg['Servers'][%[1]d]['password']  = 'password';
$cfg['Servers'][%[1]d]['verbose']   = '%[3]s';
`
		buffer.WriteString(fmt.Sprintf(config, i + 1, server["host"], server["label"]))
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

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
			Image: img.Name,
			Env: []string{
				"PMA_USER_CONFIG_BASE64=" + getBase64EncodedPMAConfig(),
				"UPLOAD_LIMIT=1024MiB",
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
