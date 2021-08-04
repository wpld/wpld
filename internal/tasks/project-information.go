package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"

	"wpld/internal/docker"
	"wpld/internal/entities"
	"wpld/internal/entities/services"
	"wpld/internal/pipelines"
	"wpld/internal/stdout"
)

func ProjectInformationPipe(api docker.Docker) pipelines.Pipe {
	getServiceInfo := func(ctx context.Context, service entities.Service, serviceName string, longestNameLength int) ([]string, []string) {
		container, err := api.ContainerInspect(ctx, service)
		if err != nil {
			return nil, nil
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
			return nil, nil
		}

		var lines, formatted []string

		for i, address := range addresses {
			if i == 0 {
				lines = append(
					lines,
					fmt.Sprintf(
						"  %s:%s %s",
						serviceName,
						strings.Repeat(" ", longestNameLength-len(serviceName)),
						address,
					),
				)

				formatted = append(
					formatted,
					fmt.Sprintf(
						"  %s:%s %s",
						serviceName,
						strings.Repeat(" ", longestNameLength-len(serviceName)),
						color.CyanString(address),
					),
				)
			} else {
				lines = append(
					lines,
					fmt.Sprintf("  %s%s", strings.Repeat(" ", longestNameLength+2), address),
				)

				formatted = append(
					formatted,
					fmt.Sprintf("  %s%s", strings.Repeat(" ", longestNameLength+2), color.CyanString(address)),
				)
			}
		}

		return lines, formatted
	}

	return func(ctx context.Context, next pipelines.NextPipe) error {
		project, ok := ctx.Value("project").(entities.Project)
		if !ok {
			return ProjectNotFoundErr
		}

		projectServices, err := project.GetServices()
		if err != nil {
			return err
		}

		lines := []string{
			fmt.Sprintf("%s project started", project.Name),
			"",
			"Project services:",
		}

		formatted := []string{
			color.New(color.FgHiWhite, color.Bold).Sprintf("%s project started", project.Name),
			"",
			"Project services:",
		}

		phpmyadminName := "phpMyAdmin"

		longest := len(phpmyadminName)
		for _, service := range projectServices {
			length := len(service.Spec.Name)
			if length > longest {
				longest = length
			}
		}

		for _, service := range projectServices {
			if service.Spec.Name != "" {
				serviceLines, serviceFormatted := getServiceInfo(ctx, service, service.Spec.Name, longest)
				if serviceLines != nil && serviceFormatted != nil {
					lines = append(lines, serviceLines...)
					formatted = append(formatted, serviceFormatted...)
				}
			}
		}

		lines = append(lines, "", "Global services:")
		formatted = append(formatted, "", "Global services:")

		phpmyadmin := services.NewPHPMyAdminService()
		serviceLines, serviceFormatted := getServiceInfo(ctx, phpmyadmin, phpmyadminName, longest)
		if serviceLines != nil && serviceFormatted != nil {
			lines = append(lines, serviceLines...)
			formatted = append(formatted, serviceFormatted...)
		}

		stdout.Box(lines, formatted)

		return next(ctx)
	}
}
