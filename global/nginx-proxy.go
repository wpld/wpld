package global

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"wpld/models"
)

const (
	NGINXPROXY_IMAGE_NAME     = "jwilder/nginx-proxy:alpine"
	NGINXPROXY_CONTAINER_NAME = "wpld_global_nginxproxy"
)

func RunNginxProxy(factory models.DockerFactory, pull bool) error {
	if pull {
		img := factory.Image(NGINXPROXY_IMAGE_NAME)
		if err := img.Pull(); err != nil {
			return err
		}
	}

	config := &container.Config{
		Image: NGINXPROXY_IMAGE_NAME,
	}

	host := &container.HostConfig{
		NetworkMode: NETWORK_NAME,
		IpcMode:     "shareable",
		Binds: []string{
			// FIXME: it won't work on Windows
			"/var/run/docker.sock:/tmp/docker.sock:ro",
		},
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "80",
				},
			},
			"443/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "443",
				},
			},
		},
	}

	nginx := factory.Container(NGINXPROXY_CONTAINER_NAME)
	if err := nginx.Create(config, host); err != nil {
		return err
	}

	return nginx.Start()
}

func StopNginxProxy(factory models.DockerFactory, rm bool) error {
	nginx := factory.Container(NGINXPROXY_CONTAINER_NAME)

	if rm {
		return nginx.Remove()
	}

	return nginx.Stop()
}
