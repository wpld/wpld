package cases

import (
	"context"

	"wpld/internal/connectors/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func ReloadProxyPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		_, err := api.FindHTTPContainers(ctx)
		if err != nil {
			return err
		}

		nginx := entities.Service{
			ID: "wpld_global__nginx_proxy",
			Spec: entities.Specification{
				Image: "jwilder/nginx-proxy:alpine",
				Ports: []string{
					"127.0.0.1:443:443",
					"127.0.0.1:80:80",
				},
				Volumes: []string{
					"/var/run/docker.sock:/tmp/docker.sock:ro",
				},
			},
		}

		if err := api.StopContainer(ctx, nginx); err != nil {
			return err
		}

		if err := api.StartContainer(ctx, nginx, false); err != nil {
			return err
		}

		return next(ctx)
	}
}
