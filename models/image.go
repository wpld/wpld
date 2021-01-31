package models

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type ImageAPI interface {
	Pull() error
}

type Image struct {
	ctx  context.Context
	cli  client.ImageAPIClient
	name string
}

func (i Image) Pull() error {
	out, err := i.cli.ImagePull(i.ctx, i.name, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	// TODO: Replace with a better approach to display image pulling progress
	if _, err = io.Copy(os.Stdout, out); err != nil {
		return err
	}

	return nil
}
