package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"wpld/internal/entities"
	"wpld/internal/misc"
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
		Attachable:     false,
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
	exists, err := d.ContainerExists(ctx, service)
	if err != nil {
		return err
	} else if exists {
		return nil
	}

	if err := d.EnsureImageExists(ctx, service.Spec.Image, pull); err != nil {
		return err
	}

	config := &container.Config{
		Cmd:          service.Spec.Cmd,
		Healthcheck:  nil,
		Image:        service.Spec.Image,
		WorkingDir:   service.Spec.WorkingDir,
		Entrypoint:   service.Spec.Entrypoint,
		AttachStderr: service.AttachStderr,
		AttachStdin:  service.AttachStdin,
		AttachStdout: service.AttachStdout,
		Labels: map[string]string{
			"wpld": misc.VERSION,
		},
	}

	if service.Project != "" {
		config.Labels["wpld.project"] = service.Project
		config.Labels["wpld.service"] = service.Spec.Name
		config.Labels["wpld.domains"] = strings.Join(service.Domains, " ")
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

	host := &container.HostConfig{
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

	var networking *network.NetworkingConfig
	if len(service.Aliases) > 0 {
		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				service.Network: {
					Aliases: service.Aliases,
				},
			},
		}
	}

	resp, err := d.api.ContainerCreate(ctx, config, host, networking, nil, service.ID)
	if err != nil {
		return err
	}

	for _, warn := range resp.Warnings {
		logrus.Warn(warn)
	}

	return nil
}

func (d Docker) ContainerExists(ctx context.Context, service entities.Service) (bool, error) {
	_, err := d.api.ContainerInspect(ctx, service.ID)
	if err == nil {
		return true, nil
	} else if !client.IsErrNotFound(err) {
		return false, err
	}

	return false, nil
}

func (d Docker) StartContainer(ctx context.Context, service entities.Service, pull bool) error {
	if err := d.EnsureContainerExists(ctx, service, pull); err != nil {
		return err
	}

	if err := d.api.ContainerStart(ctx, service.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if service.Spec.Name != "" {
		logrus.Infof("%s started", service.Spec.Name)
	}

	return nil
}

func (d Docker) StopContainer(ctx context.Context, service entities.Service) error {
	exists, err := d.ContainerExists(ctx, service)
	if err != nil {
		return err
	} else if !exists {
		return nil
	}

	if err := d.api.ContainerStop(ctx, service.ID, nil); err != nil {
		return err
	}

	if service.Spec.Name != "" {
		logrus.Infof("%s stopped", service.Spec.Name)
	}

	return nil
}

func (d Docker) StopAllContainers(ctx context.Context) error {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "wpld")

	args := types.ContainerListOptions{
		Filters: filterArgs,
	}

	list, err := d.api.ContainerList(ctx, args)
	if err != nil {
		return err
	}

	for _, c := range list {
		if err := d.api.ContainerStop(ctx, c.ID, nil); err != nil {
			return err
		}

		project, hasProjectLabel := c.Labels["wpld.project"]
		service, hasServiceLabel := c.Labels["wpld.service"]
		if hasProjectLabel && hasServiceLabel {
			logrus.Infof("{%s} %s stopped", project, service)
		}
	}

	return nil
}

func (d Docker) RestartContainer(ctx context.Context, service entities.Service) error {
	if err := d.StopContainer(ctx, service); err != nil {
		return err
	}

	for i := 0; i < 60; i++ {
		exists, err := d.ContainerExists(ctx, service)
		if err != nil {
			return err
		} else if !exists {
			break
		} else {
			time.Sleep(time.Second)
		}
	}

	if err := d.StartContainer(ctx, service, false); err != nil {
		return err
	}

	return nil
}

func (d Docker) RunContainer(ctx context.Context, service entities.Service) error {
	if err := d.StartContainer(ctx, service, false); err != nil {
		return err
	}

	statusCh, errCh := d.api.ContainerWait(ctx, service.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

func (d Docker) FindHTTPContainers(ctx context.Context) (map[string]string, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "wpld.project")
	filterArgs.Add("expose", "80")

	args := types.ContainerListOptions{
		Filters: filterArgs,
	}

	list, err := d.api.ContainerList(ctx, args)
	if err != nil {
		return nil, err
	}

	domainMapping := make(map[string]string)
	for _, c := range list {
		if domains, ok := c.Labels["wpld.domains"]; ok {
			domainMapping[domains] = c.NetworkSettings.Networks[c.HostConfig.NetworkMode].IPAddress
		}
	}

	return domainMapping, nil
}
