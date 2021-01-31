package cases

import "github.com/docker/docker/client"

func GetDockerClient() (client.CommonAPIClient, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
}
