package models

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type NetworkAPI interface {
	Inspect() (types.NetworkResource, error)
	Create() error
}

type Network struct {
	ctx  context.Context
	cli  client.NetworkAPIClient
	name string
}

func (n Network) Inspect() (types.NetworkResource, error) {
	return n.cli.NetworkInspect(n.ctx, n.name, types.NetworkInspectOptions{})
}

func (n Network) Create() error {
	inspect, err := n.Inspect()
	if err != nil {
		if !client.IsErrNotFound(err) {
			return err
		}
	} else if len(inspect.ID) > 0 {
		logrus.Debugf("[%s] Network already exists...", n.name)
		return nil
	}

	args := types.NetworkCreate{
		CheckDuplicate: false,
		Internal:       false,
		Attachable:     true,
	}

	logrus.Debugf("[%s] Creating a network...", n.name)
	if _, err = n.cli.NetworkCreate(n.ctx, n.name, args); err != nil {
		return err
	}

	return nil
}
