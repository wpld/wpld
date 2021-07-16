package tasks

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gosimple/slug"

	"wpld/internal/entities"
	"wpld/internal/entities/specs"
	"wpld/internal/pipelines"
)

func NewProjectPromptPipe() pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		var answers struct {
			Name    string
			Domains []string
			PHP     string `survey:"php"`
			Cache   string
		}

		questions := []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "Project name:"},
				Validate: survey.Required,
			},
			{
				Name:     "domains",
				Prompt:   &survey.Input{Message: "Domain names:"},
				Validate: survey.Required,
				Transform: func(answer interface{}) interface{} {
					if domains, ok := answer.(string); ok {
						return regexp.MustCompile(`[,\s]+`).Split(domains, -1)
					} else {
						return answer
					}
				},
			},
			{
				Name: "php",
				Prompt: &survey.Select{
					Message: "PHP version:",
					Default: "7.4",
					Options: []string{
						"8.0",
						"7.4",
						"7.3",
					},
				},
			},
			{
				Name: "cache",
				Prompt: &survey.Select{
					Message: "Object caching system:",
					Options: []string{
						"Memcached",
						"Redis",
						"(none)",
					},
				},
			},
		}

		err := survey.Ask(questions, &answers)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil
			} else {
				return err
			}
		}

		projectSlug := slug.Make(answers.Name)

		wpVolume := fmt.Sprintf("%s__wp", projectSlug)
		dbVolume := fmt.Sprintf("%s__db", projectSlug)

		volumes := []string{
			wpVolume,
			dbVolume,
		}

		services := map[string]entities.Specification{
			"wp":    specs.NewWordPressSpec(projectSlug, wpVolume, answers.PHP),
			"db":    specs.NewDatabaseSpec(projectSlug, dbVolume),
			"nginx": specs.NewNginxSpec(),
		}

		if answers.Cache == "Memcached" {
			services["memcache"] = specs.NewMemcachedSpec()
			if extra, ok := services["wp"].Env["WORDPRESS_CONFIG_EXTRA"]; ok {
				services["wp"].Env["WORDPRESS_CONFIG_EXTRA"] = extra + "\n$memcached_servers = array( 'memcache:11211' );"
			} else {
				services["wp"].Env["WORDPRESS_CONFIG_EXTRA"] = "$memcached_servers = array( 'memcache:11211' );"
			}
		} else if answers.Cache == "Redis" {
			services["redis"] = specs.NewRedisSpec()
			if extra, ok := services["wp"].Env["WORDPRESS_CONFIG_EXTRA"]; ok {
				services["wp"].Env["WORDPRESS_CONFIG_EXTRA"] = extra + "\ndefine( 'WP_REDIS_HOST', 'redis' );"
			} else {
				services["wp"].Env["WORDPRESS_CONFIG_EXTRA"] = "define( 'WP_REDIS_HOST', 'redis' );"
			}
		}

		return next(context.WithValue(
			ctx,
			"project",
			entities.Project{
				ID:       projectSlug,
				Name:     answers.Name,
				Domains:  answers.Domains,
				Volumes:  volumes,
				Services: services,
			},
		))
	}
}
