package tasks

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/spf13/afero"

	"wpld/internal/docker"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
)

//go:embed embeds/dnsmasq/dnsmasq.conf
var proxyConf string

func DNSReloadPipe(api docker.Docker, fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		domains, err := api.FindHTTPContainers(ctx)
		if err != nil {
			return err
		}

		dns := services.NewDnsService()
		if len(domains) == 0 {
			if err := api.ContainerStop(ctx, dns); err != nil {
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

		dns.Spec.Volumes = []string{
			fmt.Sprintf("%s:/etc/dnsmasq.conf:cached", file.Name()),
		}

		if err := api.ContainerRestart(ctx, dns); err != nil {
			return err
		}

		return next(ctx)
	}
}
