package global

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	"wpld/config"
	"wpld/models"
)

const (
	MYADMIN_IMAGE_NAME     = "phpmyadmin:latest"
	MYADMIN_CONTAINER_NAME = "wpld_global_phpmyadmin"
)

func getBase64EncodedPMAConfig() string {
	servers := []map[string]string{
		{
			"host":  MYSQL_CONTAINER_NAME,
			"label": "global",
		},
	}

	buffer := bytes.NewBufferString(`
<?php

$cfg['blowfish_secret'] = 'l3+wF5o$MUK@hj;[HLkQ4#V9-m?b4JmgXa]H_{uH#H]x|oQI%c1s|wFOGTc[<{3M';
$cfg['ServerDefault']   = 1;
`)

	for i, server := range servers {
		serverConfig := `
$cfg['Servers'][%[1]d]['host']      = '%[2]s';
$cfg['Servers'][%[1]d]['auth_type'] = 'config';
$cfg['Servers'][%[1]d]['user']      = 'root';
$cfg['Servers'][%[1]d]['password']  = 'password';
$cfg['Servers'][%[1]d]['verbose']   = '%[3]s';
`
		buffer.WriteString(fmt.Sprintf(serverConfig, i+1, server["host"], server["label"]))
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func RunMyAdmin(factory models.DockerFactory, pull bool) error {
	if pull {
		img := factory.Image(MYADMIN_IMAGE_NAME)
		if err := img.Pull(); err != nil {
			return err
		}
	}

	port := nat.PortBinding{
		HostIP:   "127.0.0.1",
		HostPort: "8092",
	}

	if viper.IsSet(config.PHPMYADMIN_PORT) {
		port.HostPort = viper.GetString(config.PHPMYADMIN_PORT)
	}

	containerConfig := &container.Config{
		Image: MYADMIN_IMAGE_NAME,
		Env: []string{
			"PMA_USER_CONFIG_BASE64=" + getBase64EncodedPMAConfig(),
			"UPLOAD_LIMIT=" + viper.GetString(config.PHPMYADMIN_UPLOAD_LIMIT),
		},
	}

	host := &container.HostConfig{
		NetworkMode: NETWORK_NAME,
		IpcMode:     "shareable",
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{port},
		},
	}

	myadmin := factory.Container(MYADMIN_CONTAINER_NAME)
	if err := myadmin.Create(containerConfig, host); err != nil {
		return err
	}

	return myadmin.Start()
}

func StopMyAdmin(factory models.DockerFactory, rm bool) error {
	myadmin := factory.Container(MYADMIN_CONTAINER_NAME)

	if rm {
		return myadmin.Remove()
	}

	return myadmin.Stop()
}
