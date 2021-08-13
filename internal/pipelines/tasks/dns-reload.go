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
	"wpld/internal/stdout"
)

//go:embed embeds/dnsmasq/dnsmasq.conf
var dnsmasqConf string

func DNSReloadPipe(api docker.Docker, fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		domains, err := api.FindContainersWithDomains(ctx)
		if err != nil {
			return err
		}

		dns := services.NewDnsService()
		if len(domains) == 0 {
			stdout.StartSpinner("Stopping global DNS...")
			err := api.ContainerStop(ctx, dns)
			stdout.StopSpinner()

			if err != nil {
				return err
			} else {
				stdout.Success("Global DNS stopped")
			}

			return next(ctx)
		}

		tmpdir := afero.GetTempDir(fs, "wpld")
		file, err := afero.TempFile(fs, tmpdir, "dnsmasq.*.conf")
		if err != nil {
			return err
		}

		tmpl, err := template.New("dnsmasq").Parse(dnsmasqConf)
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

		stdout.StartSpinner("Starting global DNS...")
		err = api.ContainerRestart(ctx, dns)
		stdout.StopSpinner()

		if err != nil {
			return err
		} else {
			stdout.Success("Global DNS started")
		}

		return next(ctx)
	}
}
