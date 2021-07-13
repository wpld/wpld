package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"wpld/internal/entities"
)

type Docker struct {
	api client.CommonAPIClient
}

func NewDocker() (Docker, error) {
	api, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return Docker{}, err
	}

	docker := Docker{
		api: api,
	}

	return docker, nil
}

func (d Docker) EnsureImageExists(ctx context.Context, imageID string, force bool) error {
	if !force {
		images, err := d.api.ImageList(ctx, types.ImageListOptions{})
		if err != nil {
			return err
		}

		for _, image := range images {
			for _, tag := range image.RepoTags {
				if tag == imageID {
					return nil
				}
			}
		}
	}

	out, err := d.api.ImagePull(ctx, imageID, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, _ = io.Copy(os.Stdout, out)

	return nil
}

func (d Docker) EnsureNetworkExists(ctx context.Context, networkID string) error {
	_, err := d.api.NetworkInspect(ctx, networkID, types.NetworkInspectOptions{})
	if err == nil {
		return nil
	} else if !client.IsErrNotFound(err) {
		return err
	}

	args := types.NetworkCreate{
		CheckDuplicate: false,
		Driver:         "bridge",
		EnableIPv6:     false,
		Internal:       false,
		Attachable:     true,
		Labels:         basicLabels,
	}

	_, err = d.api.NetworkCreate(ctx, networkID, args)

	return err
}

func (d Docker) EnsureVolumeExists(ctx context.Context, volumeID string) error {
	_, err := d.api.VolumeInspect(ctx, volumeID)
	if err == nil {
		return nil
	} else if !client.IsErrNotFound(err) {
		return err
	}

	args := volume.VolumeCreateBody{
		Name:   volumeID,
		Driver: "local",
		Labels: basicLabels,
	}

	_, err = d.api.VolumeCreate(ctx, args)

	return err
}

func (d Docker) EnsureContainerExists(ctx context.Context, service entities.Service, pull bool) error {
	_, err := d.api.ContainerInspect(ctx, service.ID)
	if err == nil {
		return nil
	} else if !client.IsErrNotFound(err) {
		return err
	}

	if err := d.EnsureImageExists(ctx, service.Spec.Image, pull); err != nil {
		return err
	}

	config := container.Config{
		Cmd:         service.Spec.Cmd,
		Healthcheck: nil,
		Image:       service.Spec.Image,
		WorkingDir:  service.Spec.WorkingDir,
		Entrypoint:  service.Spec.Entrypoint,
	}

	envLen := len(service.Spec.Env)
	if envLen > 0 {
		i := 0
		config.Env = make([]string, envLen)
		for key, value := range service.Spec.Env {
			config.Env[i] = fmt.Sprintf("%s=%s", key, value)
			i++
		}
	}

	host := container.HostConfig{
		Binds:       NormalizeContainerBinds(service.Spec.Volumes),
		NetworkMode: container.NetworkMode(service.Network),
		AutoRemove:  true,
		IpcMode:     "shareable",
		CapAdd:      service.Spec.CapAdd,
		CapDrop:     service.Spec.CapDrop,
		VolumesFrom: service.Spec.VolumesFrom,
	}

	if len(service.Spec.Ports) > 0 {
		_, portBindings, err := nat.ParsePortSpecs(service.Spec.Ports)
		if err != nil {
			return err
		} else {
			host.PortBindings = portBindings
		}
	}

	resp, err := d.api.ContainerCreate(ctx, &config, &host, nil, nil, service.ID)
	if err != nil {
		return err
	}

	for _, warn := range resp.Warnings {
		logrus.Warn(warn)
	}

	return nil
}

func (d Docker) StartContainer(ctx context.Context, service entities.Service, pull bool) error {
	if err := d.EnsureContainerExists(ctx, service, pull); err != nil {
		return err
	}

	if err := d.api.ContainerStart(ctx, service.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logrus.Infof("%s started", service.Spec.Name)

	return nil
}

func (d Docker) StopContainer(ctx context.Context, service entities.Service) error {
	_, err := d.api.ContainerInspect(ctx, service.ID)
	if client.IsErrNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	if err := d.api.ContainerStop(ctx, service.ID, nil); err != nil {
		return err
	}

	logrus.Infof("%s stopped", service.Spec.Name)

	return nil
}
