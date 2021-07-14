package cases

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gosimple/slug"

	"wpld/internal/entities"
	"wpld/internal/pipelines"
)

func NewProjectPromptPipe() pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		var answers struct {
			Name    string
			Domains []string
			PHP     string `survey:"php"`
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

		wp := entities.Specification{
			Name:  "WordPress",
			Image: fmt.Sprintf("wordpress:5-php%s-fpm-alpine", answers.PHP),
			Volumes: []string{
				fmt.Sprintf("%s:/var/www/html", wpVolume),
			},
			Env: map[string]string{
				"WORDPRESS_DB_HOST":     "db",
				"WORDPRESS_DB_USER":     "wordpress",
				"WORDPRESS_DB_PASSWORD": "password",
				"WORDPRESS_DB_NAME":     projectSlug,
			},
		}

		db := entities.Specification{
			Name:  "Database",
			Image: "mariadb:latest",
			Volumes: []string{
				fmt.Sprintf("%s:/var/lib/mysql", dbVolume),
			},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "password",
				"MYSQL_DATABASE":      projectSlug,
				"MYSQL_USER":          "wordpress",
				"MYSQL_PASSWORD":      "password",
			},
		}

		nginx := entities.Specification{
			Name:    "Nginx",
			Image:   "nginx:alpine",
			Volumes: []string{},
			VolumesFrom: []string{
				"wp",
			},
			DependsOn: []string{
				"wp",
			},
		}

		return next(context.WithValue(
			ctx,
			"project",
			entities.Project{
				ID:      projectSlug,
				Name:    answers.Name,
				Domains: answers.Domains,
				Volumes: []string{
					wpVolume,
					dbVolume,
				},
				Services: map[string]entities.Specification{
					"wp":    wp,
					"db":    db,
					"nginx": nginx,
				},
			},
		))
	}
}
