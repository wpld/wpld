package tasks

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"
	"wpld/internal/stdout"

	"github.com/spf13/afero"

	"wpld/internal/docker"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

//go:embed embeds/nginx/reverse-proxy.conf
var proxyConf string

func GlobalProxyReload(api docker.Docker, fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		domains, err := api.FindContainersWithDomains(ctx)
		if err != nil {
			return err
		}

		proxy := services.NewProxyService()
		if len(domains) == 0 {
			if err := api.ContainerStop(ctx, proxy); err != nil {
				return err
			} else {
				return next(ctx)
			}
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

		stdout.StartSpinner("Starting global proxy...")
		err = api.ContainerRestart(ctx, proxy)
		if err == nil {
			if networks, networksErr := api.FindAllNetworks(ctx); networksErr == nil {
				err = api.ContainerConnectNetworks(ctx, proxy, networks)
			}
		}
		stdout.StopSpinner()

		if err != nil {
			return err
		} else {
			stdout.Success("Global proxy started")
		}

		return next(ctx)
	}
}