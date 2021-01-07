package global

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"wpld/utils"
)

const (
	NGINXPROXY_IMAGE_NAME = "jwilder/nginx-proxy:alpine"
	NGINXPROXY_CONTAINER_NAME = "wpld_global_nginxproxy"
)

func RunNginxProxy(ctx context.Context, cli *client.Client, pull bool) error {
	img := utils.Image{
		Name: NGINXPROXY_IMAGE_NAME,
	}

	if pull {
		if err := img.Pull(ctx, cli); err != nil {
			return err
		}
	}

	nginx := utils.Container{
		Name: NGINXPROXY_CONTAINER_NAME,
		Create: &container.Config{
			Image: img.Name,
		},
		Host: &container.HostConfig{
			NetworkMode: NETWORK_NAME,
			IpcMode: "shareable",
			Binds: []string{
				"/var/run/docker.sock:/tmp/docker.sock:ro",
			},
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP: "127.0.0.1",
						HostPort: "80",
					},
				},
				"443/tcp": []nat.PortBinding{
					{
						HostIP: "127.0.0.1",
						HostPort: "443",
					},
				},
			},
		},
	}

	return nginx.Start(ctx, cli)
}

func StopNginxProxy(ctx context.Context, cli *client.Client, rm bool) error {
	nginx := utils.Container{
		Name: NGINXPROXY_CONTAINER_NAME,
	}

	if rm {
		return nginx.Remove(ctx, cli)
	}

	return nginx.Stop(ctx, cli)
}
