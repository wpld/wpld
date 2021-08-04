package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func ProjectInformationPipe(api docker.Docker) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		services, err := project.GetServices()
		if err != nil {
			return err
		}

		lines := []string{
			fmt.Sprintf("%s is started", project.Name),
			"",
		}

		formatted := []string{
			color.New(color.FgHiWhite, color.Bold).Sprintf("%s is started", project.Name),
			"",
		}

		longest := 0
		for _, service := range services {
			length := len(service.Spec.Name)
			if length > longest {
				longest = length
			}
		}

		for _, service := range services {
			if service.Spec.Name == "" {
				continue
			}

			container, err := api.ContainerInspect(ctx, service)
			if err != nil {
				continue
			}

			ip := ""
			if networkInfo, ok := container.NetworkSettings.Networks[string(container.HostConfig.NetworkMode)]; ok {
				ip = networkInfo.IPAddress
			}

			var addresses []string

			if ip != "" {
				for rawPort := range container.Config.ExposedPorts {
					proto, port := nat.SplitProtoPort(string(rawPort))
					if proto == "tcp" {
						if port != "80" {
							addresses = append(addresses, fmt.Sprintf("%s:%s", ip, port))
						} else {
							addresses = append(addresses, ip)
						}
					}
				}
			}

			if service.Spec.Domains != nil {
				for _, domain := range service.Spec.Domains {
					for rawPort := range container.Config.ExposedPorts {
						proto, port := nat.SplitProtoPort(string(rawPort))
						if proto == "tcp" {
							if port != "80" {
								addresses = append(addresses, fmt.Sprintf("%s:%s", domain, port))
							} else {
								addresses = append(addresses, domain)
							}
						}
					}
				}
			}

			if len(addresses) == 0 {
				continue
			}

			for i, address := range addresses {
				if i == 0 {
					lines = append(
						lines,
						fmt.Sprintf(
							"%s:%s %s",
							service.Spec.Name,
							strings.Repeat(" ", longest-len(service.Spec.Name)),
							address,
						),
					)

					formatted = append(
						formatted,
						fmt.Sprintf(
							"%s:%s %s",
							service.Spec.Name,
							strings.Repeat(" ", longest-len(service.Spec.Name)),
							color.CyanString(address),
						),
					)
				} else {
					lines = append(
						lines,
						strings.Repeat(" ", longest+2)+address,
					)

					formatted = append(
						formatted,
						strings.Repeat(" ", longest+2)+color.CyanString(address),
					)
				}
			}
		}

		stdout.Box(lines, formatted)

		return next(ctx)
	}
}
