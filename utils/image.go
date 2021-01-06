package utils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type Image struct {
	Name string
	PullOptions types.ImagePullOptions
}

func (args Image) Pull(ctx context.Context, cli *client.Client) error {
	if out, err := cli.ImagePull(ctx, args.Name, args.PullOptions); err != nil {
		return err
	} else {
		// TODO: Replace with a better approach to display image pulling progress
		io.Copy(os.Stdout, out)
	}

	return nil
}
