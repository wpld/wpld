package global

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"wpld/models"
)

const (
	DNSMASQ_IMAGE_NAME     = "andyshinn/dnsmasq:latest"
	DNSMASQ_CONTAINER_NAME = "wpld_global_dnsmasq"
)

func RunDnsMasq(factory models.DockerFactory, pull bool) error {
	if pull {
		img := factory.Image(DNSMASQ_IMAGE_NAME)
		if err := img.Pull(); err != nil {
			return err
		}
	}

	portBinding := []nat.PortBinding{
		{
			HostIP:   "127.0.0.1",
			HostPort: "53",
		},
	}

	nginx := factory.Container(NGINXPROXY_CONTAINER_NAME)
	nginxContainer, nginxErr := nginx.Inspect()
	if nginxErr != nil {
		return nginxErr
	}

	config := &container.Config{
		Image: DNSMASQ_IMAGE_NAME,
		Cmd: []string{
			"-A",
			"/test/" + nginxContainer.NetworkSettings.Networks[NETWORK_NAME].IPAddress,
			"--log-facility=-",
		},
	}

	host := &container.HostConfig{
		NetworkMode: NETWORK_NAME,
		IpcMode:     "shareable",
		CapAdd: []string{
			"NET_ADMIN",
		},
		PortBindings: nat.PortMap{
			"53/tcp": portBinding,
			"53/udp": portBinding,
		},
	}

	dnsmasq := factory.Container(DNSMASQ_CONTAINER_NAME)
	if err := dnsmasq.Create(config, host); err != nil {
		return err
	}

	return dnsmasq.Start()
}

func StopDnsMasq(factory models.DockerFactory, rm bool) error {
	dnsmasq := factory.Container(DNSMASQ_CONTAINER_NAME)

	if rm {
		return dnsmasq.Remove()
	}

	return dnsmasq.Stop()
}
