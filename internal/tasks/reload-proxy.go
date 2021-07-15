package tasks

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"
	"time"

	"github.com/spf13/afero"

	"wpld/internal/connectors/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

//go:embed embeds/nginx/reverse-proxy.conf
var proxyConf string

func ReloadProxyPipe(api docker.Docker, fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		proxy := entities.Service{
			ID:      "wpld__reverse_proxy",
			Network: "host",
			Spec: entities.Specification{
				Image: "nginx:alpine",
			},
		}

		if err := api.StopContainer(ctx, proxy); err != nil {
			return err
		}

		tmpdir := afero.GetTempDir(fs, "wpld")
		file, err := afero.TempFile(fs, tmpdir, "reverse-proxy.*.conf")
		if err != nil {
			return err
		}

		tmpl, err := template.New("proxy").Parse(proxyConf)
		if err != nil {
			return err
		}

		domains, err := api.FindHTTPContainers(ctx)
		if err != nil {
			return err
		}

		buff := bytes.NewBufferString("")
		if err := tmpl.Execute(buff, domains); err != nil {
			return err
		}

		if _, err := file.Write(buff.Bytes()); err != nil {
			return err
		}

		proxy.Spec.Volumes = []string{
			fmt.Sprintf("%s:/etc/nginx/conf.d/default.conf:cached", file.Name()),
		}

		for i := 0; i < 60; i++ {
			exists, err := api.ContainerExists(ctx, proxy)
			if err != nil {
				return err
			} else if !exists {
				break
			} else {
				time.Sleep(time.Second)
			}
		}

		if err := api.StartContainer(ctx, proxy, false); err != nil {
			return err
		}

		return next(ctx)
	}
}
