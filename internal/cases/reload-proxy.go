package cases

import (
	"context"

	"wpld/internal/connectors/docker"
	"wpld/internal/pipelines"
)

func ReloadProxyPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		err := api.FindHTTPContainers(ctx)
		if err != nil {
			return err
		}

		return next(ctx)
	}
}
