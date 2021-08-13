package docker

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"wpld/internal/entities"
	"wpld/internal/stdout"
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

func (d Docker) EnsureNetworkExists(ctx context.Context, net entities.Network) error {
	_, err := d.api.NetworkInspect(ctx, net.Name, types.NetworkInspectOptions{})
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
		Labels:         GetBasicLabels(),
	}

	if net.Subnet != "" {
		args.IPAM = &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet: net.Subnet,
				},
			},
		}
	}

	_, err = d.api.NetworkCreate(ctx, net.Name, args)

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
		Labels: GetBasicLabels(),
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

	if service.Network.Name != "" {
		if err := d.EnsureNetworkExists(ctx, service.Network); err != nil {
			return err
		}
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
		OpenStdin:    service.AttachStdin,
		Tty:          service.Tty,
		Labels:       GetBasicLabels(),
	}

	exposedPortsLen := len(service.Spec.ExposedPorts)
	if exposedPortsLen > 0 {
		config.ExposedPorts = make(map[nat.Port]struct{}, exposedPortsLen)

		for _, exposedPort := range service.Spec.ExposedPorts {
			proto, port := nat.SplitProtoPort(exposedPort)
			start, end, err := nat.ParsePortRangeToInt(port)
			if err != nil {
				return err
			}

			for i := start; i <= end; i++ {
				p, err := nat.NewPort(proto, strconv.Itoa(i))
				if err != nil {
					return err
				}

				if _, exists := config.ExposedPorts[p]; !exists {
					config.ExposedPorts[p] = struct{}{}
				}
			}
		}
	}

	if service.Project != "" {
		config.Labels["io.wpld.project"] = service.Project
	}

	if service.Spec.Name != "" {
		config.Labels["io.wpld.service"] = service.Spec.Name
	}

	if len(service.Spec.Domains) > 0 {
		config.Labels["io.wpld.domains"] = strings.Join(service.Spec.Domains, ",")
	}

	envLen := len(service.Spec.Env)
	if envLen > 0 {
		config.Env = service.Spec.GetEnvs()
	}

	host := &container.HostConfig{
		Binds:       NormalizeContainerBinds(service.Spec.Volumes),
		NetworkMode: container.NetworkMode(service.Network.Name),
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

	if len(service.Aliases) > 0 || service.Spec.IPAddress != "" {
		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				service.Network.Name: {
					Aliases: service.Aliases,
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: service.Spec.IPAddress,
					},
				},
			},
		}
	}

	resp, err := d.api.ContainerCreate(ctx, config, host, networking, nil, service.ID)
	if err != nil {
		return err
	}

	for _, warn := range resp.Warnings {
		stdout.Warn(warn)
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

func (d Docker) ContainerStart(ctx context.Context, service entities.Service, pull bool) error {
	if err := d.EnsureContainerExists(ctx, service, pull); err != nil {
		return err
	}

	if err := d.api.ContainerStart(ctx, service.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerStop(ctx context.Context, service entities.Service) error {
	exists, err := d.ContainerExists(ctx, service)
	if err != nil {
		return err
	} else if !exists {
		return nil
	}

	if err := d.api.ContainerStop(ctx, service.ID, nil); err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerRestart(ctx context.Context, service entities.Service) error {
	if err := d.ContainerStop(ctx, service); err != nil {
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

	if err := d.ContainerStart(ctx, service, false); err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerAttach(ctx context.Context, service entities.Service) error {
	attachOptions := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  service.AttachStdin,
		Stdout: service.AttachStdout,
		Stderr: service.AttachStderr,
	}

	attach, err := d.api.ContainerAttach(ctx, service.ID, attachOptions)
	if err != nil {
		return err
	}

	// @see: https://github.com/docker/cli/blob/master/cli/command/container/attach.go

	// TODO: add signals forwarding https://github.com/docker/cli/blob/master/cli/command/container/attach.go#L99-L102

	// TODO: fix broken raw terminal state

	hijackedStreamer := NewHijackedStreamer(attach, service.Tty)
	if err := hijackedStreamer.Stream(ctx); err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerExecAttach(ctx context.Context, service entities.Service, cmd []string, wd string) error {
	execOptions := types.ExecConfig{
		Tty:          service.Tty,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Env:          service.Spec.GetEnvs(),
		WorkingDir:   wd,
		Cmd:          cmd,
	}

	idresp, err := d.api.ContainerExecCreate(ctx, service.ID, execOptions)
	if err != nil {
		return err
	}

	attachOptions := types.ExecStartCheck{
		Tty: service.Tty,
	}

	attach, err := d.api.ContainerExecAttach(ctx, idresp.ID, attachOptions)
	if err != nil {
		return err
	}

	hijackedStreamer := NewHijackedStreamer(attach, service.Tty)
	if err := hijackedStreamer.Stream(ctx); err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerLogs(ctx context.Context, service entities.Service, tail string, skipStdout, skipStderr bool) error {
	params := types.ContainerLogsOptions{
		ShowStdout: !skipStdout,
		ShowStderr: !skipStderr,
		Tail:       tail,
	}

	resp, err := d.api.ContainerLogs(ctx, service.ID, params)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return err
	}

	return nil
}

func (d Docker) ContainerInspect(ctx context.Context, service entities.Service) (types.ContainerJSON, error) {
	return d.api.ContainerInspect(ctx, service.ID)
}

func (d Docker) ContainerConnectNetworks(ctx context.Context, service entities.Service, networks []string) error {
	for _, n := range networks {
		config := network.EndpointSettings{}
		if err := d.api.NetworkConnect(ctx, n, service.ID, &config); err != nil {
			return err
		}
	}

	return nil
}

func (d Docker) NetworkIsInUsed(ctx context.Context, net string) (bool, error) {
	info, err := d.api.NetworkInspect(ctx, net, types.NetworkInspectOptions{})
	if err != nil {
		return false, err
	} else {
		return len(info.Containers) > 0, nil
	}
}

func (d Docker) NetworkRemove(ctx context.Context, net string) error {
	return d.api.NetworkRemove(ctx, net)
}

func (d Docker) FindAllRunningContainers(ctx context.Context) ([]types.Container, error) {
	args := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg(
				"label",
				"io.wpld.version",
			),
		),
	}

	return d.api.ContainerList(ctx, args)
}

func (d Docker) FindMySQLContainers(ctx context.Context) (map[string]string, error) {
	args := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg(
				"label",
				"io.wpld.version",
			),
			filters.Arg(
				"expose",
				"3306",
			),
		),
	}

	list, err := d.api.ContainerList(ctx, args)
	if err != nil {
		return nil, err
	}

	domainMapping := make(map[string]string, len(list))
	for _, c := range list {
		ip := c.NetworkSettings.Networks[c.HostConfig.NetworkMode].IPAddress

		if project, ok := c.Labels["io.wpld.project"]; ok {
			domainMapping[ip] = project
		} else {
			domainMapping[ip] = ip
		}
	}

	return domainMapping, nil
}

func (d Docker) FindContainersWithDomains(ctx context.Context) (map[string]string, error) {
	args := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg(
				"label",
				"io.wpld.domains",
			),
		),
	}

	list, err := d.api.ContainerList(ctx, args)
	if err != nil {
		return nil, err
	}

	domainMapping := make(map[string]string)
	for _, c := range list {
		if domainsLabel, ok := c.Labels["io.wpld.domains"]; ok {
			domains := strings.Split(domainsLabel, ",")
			for _, domain := range domains {
				if networkInfo, ok := c.NetworkSettings.Networks[c.HostConfig.NetworkMode]; ok {
					domainMapping[domain] = networkInfo.IPAddress
				}
			}
		}
	}

	return domainMapping, nil
}

func (d Docker) FindAllNetworks(ctx context.Context) ([]string, error) {
	args := types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg(
				"label",
				"io.wpld.version",
			),
		),
	}

	list, err := d.api.NetworkList(ctx, args)
	if err != nil {
		return nil, err
	}

	networks := []string{}
	for _, net := range list {
		if net.Name != "wpld" {
			networks = append(networks, net.ID)
		}
	}

	return networks, nil
}
