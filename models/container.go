package models

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type ContainerAPI interface {
	Create(config *container.Config, host *container.HostConfig) error
	Inspect() (types.ContainerJSON, error)
	Remove() error
	Start() error
	Stop() error
}

type Container struct {
	ctx  context.Context
	cli  client.ContainerAPIClient
	name string
}

func (c Container) Inspect() (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(c.ctx, c.name)
}

func (c Container) Create(config *container.Config, host *container.HostConfig) error {
	inspect, err := c.Inspect()
	if err != nil {
		if !client.IsErrNotFound(err) {
			return err
		}
	}

	if len(inspect.ID) > 0 {
		logrus.Debugf("[%s] Container already exists...", c.name)
		return nil
	}

	logrus.Debugf("[%s] Creating container...", c.name)
	_, err = c.cli.ContainerCreate(c.ctx, config, host, nil, nil, c.name)
	if err != nil {
		return err
	}

	return nil
}

func (c Container) Start() error {
	inspect, err := c.Inspect()
	if err != nil {
		return err
	}

	if inspect.State.Running {
		logrus.Debugf("[%s] Container is already running...", c.name)
		return nil
	}

	logrus.Debugf("[%s] Starting container...", c.name)
	return c.cli.ContainerStart(c.ctx, inspect.ID, types.ContainerStartOptions{})
}

func (c Container) Stop() error {
	inspect, err := c.Inspect()
	if err != nil {
		if client.IsErrNotFound(err) {
			logrus.Debugf("[%s] Container doesn't exist...", c.name)
			return nil
		}

		return err
	}

	logrus.Debugf("[%s] Stopping container: %s", c.name, inspect.ID)
	return c.cli.ContainerStop(c.ctx, inspect.ID, nil)
}

func (c Container) Remove() error {
	inspect, err := c.Inspect()
	if err != nil {
		if client.IsErrNotFound(err) {
			logrus.Debugf("[%s] Container doesn't exist...", c.name)
			return nil
		}

		return err
	}

	logrus.Debugf("[%s] Deleting container: %s", c.name, inspect.ID)
	return c.cli.ContainerRemove(c.ctx, inspect.ID, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
}
