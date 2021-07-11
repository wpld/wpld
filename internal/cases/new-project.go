package cases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gosimple/slug"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"wpld/internal/controllers/pipelines"
	"wpld/internal/entities"
)

func NewProjectPipeline(fs afero.Fs) pipelines.Pipeline {
	return pipelines.NewPipeline(
		newProjectPromptPipe(),
		newProjectMarshalPipe(fs),
	)
}

func newProjectPromptPipe() pipelines.Pipe {
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

		return next(context.WithValue(
			ctx,
			"project",
			entities.Project{
				Name:    answers.Name,
				Domains: answers.Domains,
				Volumes: []string{
					wpVolume,
				},
				Services: map[string]entities.Service{
					"wp": {
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
					},
					"db": {
						Name:  "Database",
						Image: "mariadb:latest",
						Ports: []string{
							"3306:3306",
						},
						Env: map[string]string{
							"MYSQL_DATABASE":           projectSlug,
							"MYSQL_USER":               "wordpress",
							"MYSQL_PASSWORD":           "password",
							"MYSQL_INITDB_SKIP_TZINFO": "skip",
						},
					},
				},
			},
		))
	}
}

func newProjectMarshalPipe(fs afero.Fs) pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		buffer := bytes.NewBufferString("")

		encoder := yaml.NewEncoder(buffer)
		encoder.SetIndent(2)

		if err := encoder.Encode(ctx.Value("project")); err != nil {
			return err
		}

		if err := afero.WriteFile(fs, ".wpld.yml", buffer.Bytes(), 0644); err != nil {
			return err
		}

		return next(ctx)
	}
}
