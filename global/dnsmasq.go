package global

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"wpld/utils"
)

const (
	DNSMASQ_IMAGE_NAME = "andyshinn/dnsmasq:latest"
	DNSMASQ_CONTAINER_NAME = "wpld_global_dnsmasq"
)

func RunDnsMasq(ctx context.Context, cli *client.Client, pull bool) error {
	img := utils.Image{
		Name: DNSMASQ_IMAGE_NAME,
	}

	if pull {
		if err := img.Pull(ctx, cli); err != nil {
			return err
		}
	}

	portBinding := []nat.PortBinding{
		{
			HostIP: "127.0.0.1",
			HostPort: "53",
		},
	}

	nginx := utils.Container{
		Name: NGINXPROXY_CONTAINER_NAME,
	}

	nginxContainer, nginxErr := nginx.Inspect(ctx, cli)
	if nginxErr != nil {
		return nginxErr
	}

	dnsmasq := utils.Container{
		Name: DNSMASQ_CONTAINER_NAME,
		Create: &container.Config{
			Image: img.Name,
			Cmd: []string{
				"-A",
				"/test/" + nginxContainer.NetworkSettings.Networks[NETWORK_NAME].IPAddress,
				"--log-facility=-",
			},
		},
		Host: &container.HostConfig{
			NetworkMode: NETWORK_NAME,
			IpcMode: "shareable",
			CapAdd: []string{
				"NET_ADMIN",
			},
			PortBindings: nat.PortMap{
				"53/tcp": portBinding,
				"53/udp": portBinding,
			},
		},
	}

	return dnsmasq.Start(ctx, cli)
}

func StopDnsMasq(ctx context.Context, cli *client.Client, rm bool) error {
	dnsmasq := utils.Container{
		Name: DNSMASQ_CONTAINER_NAME,
	}

	if rm {
		return dnsmasq.Remove(ctx, cli)
	}

	return dnsmasq.Stop(ctx, cli)
}
