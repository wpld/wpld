package models

import (
	"context"
	"github.com/docker/docker/client"
)

type DockerFactory struct {
	ctx context.Context
	cli client.CommonAPIClient
}

func NewDockerFactory(ctx context.Context, cli client.CommonAPIClient) DockerFactory {
	return DockerFactory{
		ctx: ctx,
		cli: cli,
	}
}

func (f DockerFactory) Container(name string) ContainerAPI {
	return Container{
		ctx:  f.ctx,
		cli:  f.cli,
		name: name,
	}
}

func (f DockerFactory) Image(name string) ImageAPI {
	return Image{
		ctx:  f.ctx,
		cli:  f.cli,
		name: name,
	}
}

func (f DockerFactory) Network(name string) NetworkAPI {
	return Network{
		ctx:  f.ctx,
		cli:  f.cli,
		name: name,
	}
}
