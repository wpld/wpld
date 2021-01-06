package utils

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

type Container struct {
	Name string
	Create *container.Config
	Host *container.HostConfig
	Network *network.NetworkingConfig
	Platform *v1.Platform
	StartOptions types.ContainerStartOptions
	RemoveOptions types.ContainerRemoveOptions
}

func (args Container) Inspect(ctx context.Context, cli *client.Client) (types.ContainerJSON, error) {
	logrus.Debugf("Checking {%s} container...", args.Name)
	return cli.ContainerInspect(ctx, args.Name)
}

func (args Container) Start(ctx context.Context, cli *client.Client) error {
	var containerID string
	if c, cerr := args.Inspect(ctx, cli); cerr != nil {
		logrus.Debugf("Container {%s} doesn't exist, creating...", args.Name)
		if resp, err := cli.ContainerCreate(ctx, args.Create, args.Host, args.Network, args.Platform, args.Name); err != nil {
			return err
		} else {
			logrus.Debugf("Container {%s} has been created: %s", args.Name, resp.ID)
			containerID = resp.ID
		}
	} else {
		if c.State.Running {
			logrus.Debugf("Container {%s} is up and running...", args.Name)
			return nil
		}

		if c.State.Dead {
			logrus.Debugf("Container {%s} is dead...", args.Name)
		} else {
			logrus.Debugf("Container {%s} exists: %s", args.Name, c.ID)
		}

		containerID = c.ID
	}

	logrus.Debugf("Starting {%s} container...", args.Name)
	err := cli.ContainerStart(ctx, containerID, args.StartOptions)
	if err != nil {
		return err
	}

	return nil
}

func (args Container) Stop(ctx context.Context, cli *client.Client) error {
	c, cerr := args.Inspect(ctx, cli)
	if cerr != nil {
		fmt.Printf("%v\n", client.IsErrNotFound(cerr));
		logrus.Debugf("Container {%s} doesn't exist...", args.Name)
		return nil
	}

	if c.State.Running {
		logrus.Debugf("Container {%s} is running, stopping...", args.Name)
		if serr := cli.ContainerStop(ctx, c.ID, nil); serr != nil {
			return serr
		}
	} else {
		logrus.Debugf("Container {%s} isn't running: %s", args.Name, c.ID)
	}

	return nil
}

func (args Container) Remove(ctx context.Context, cli *client.Client) error {
	c, cerr := args.Inspect(ctx, cli)
	if cerr != nil {
		fmt.Printf("%v\n", client.IsErrNotFound(cerr));
		logrus.Debugf("Container {%s} doesn't exist...", args.Name)
		return nil
	}

	if c.State.Running {
		logrus.Debugf("Container {%s} is running, stopping...", args.Name)
		if serr := cli.ContainerStop(ctx, c.ID, nil); serr != nil {
			return serr
		}
	}

	logrus.Debugf("Deleting {%s} container: %s", args.Name, c.ID)
	return cli.ContainerRemove(ctx, c.ID, args.RemoveOptions)
}
