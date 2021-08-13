package tasks

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"text/template"

	"wpld/internal/docker"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

//go:embed embeds/phpmyadmin/config.php
var phpMyAdminConf string

func PHPMyAdminReloadPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		containers, err := api.FindMySQLContainers(ctx)
		if err != nil {
			return err
		}

		phpmyadmin := services.NewPHPMyAdminService()
		if len(containers) == 0 {
			stdout.StartSpinner("Stopping global phpMyAdmin...")
			err := api.ContainerStop(ctx, phpmyadmin)
			stdout.StopSpinner()

			if err != nil {
				return err
			} else {
				stdout.Success("Global phpMyAdmin stopped")
			}

			return next(ctx)
		}

		tmpl, err := template.New("phpmyadmin").Parse(phpMyAdminConf)
		if err != nil {
			return err
		}

		buff := bytes.NewBufferString("")
		if err := tmpl.Execute(buff, containers); err != nil {
			return err
		}

		if phpmyadmin.Spec.Env == nil {
			phpmyadmin.Spec.Env = make(map[string]string, 1)
		}

		phpmyadmin.Spec.Env["PMA_USER_CONFIG_BASE64"] = base64.StdEncoding.EncodeToString(buff.Bytes())

		stdout.StartSpinner("Starting global phpMyAdmin...")
		err = api.ContainerRestart(ctx, phpmyadmin)
		if err == nil {
			if networks, networksErr := api.FindAllNetworks(ctx); networksErr == nil {
				err = api.ContainerConnectNetworks(ctx, phpmyadmin, networks)
			}
		}
		stdout.StopSpinner()

		if err != nil {
			return err
		} else {
			stdout.Success("Global phpMyAdmin started")
		}

		return next(ctx)
	}
}
