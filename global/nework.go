package global

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

const (
	NETWORK_NAME = "wpld"
)

func VerifyNetwork(ctx context.Context, cli *client.Client) (string, error) {
	logrus.Debugf("Checking {%s} network...", NETWORK_NAME)

	n, nerr := cli.NetworkInspect(ctx, NETWORK_NAME, types.NetworkInspectOptions{})
	if nerr != nil {
		createArgs := types.NetworkCreate{
			CheckDuplicate: false,
			Internal: false,
			Attachable: false,
		}

		logrus.Debugf("Creating a new {%s} network...", NETWORK_NAME)
		if r, rerr := cli.NetworkCreate(ctx, NETWORK_NAME, createArgs); rerr != nil {
			return "", rerr
		} else {
			logrus.Debugf("Network {%s} is created: %s\n", NETWORK_NAME, r.ID)
			return r.ID, nil
		}
	}

	logrus.Debugf("Network {%s} exists: %s\n", NETWORK_NAME, n.ID)
	return n.ID, nil
}
