package cases

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"

	"wpld/internal/controllers/pipelines"
	"wpld/internal/entities"
)

var NewProject = pipelines.NewPipeline(
	newProjectPromptPipe,
	newProjectMarshalPipe,
)

func newProjectPromptPipe(ctx context.Context, next pipelines.NextPipe) error {
	var answers struct {
		Name    string
		Domains []string
	}

	questions := []*survey.Question{
		{
			Name:     "Name",
			Prompt:   &survey.Input{Message: "Project name:"},
			Validate: survey.Required,
		},
		{
			Name:     "Domains",
			Prompt:   &survey.Input{Message: "Domain names:"},
			Validate: survey.Required,
			Transform: func(answer interface{}) interface{} {
				domains, ok := answer.(string)
				if !ok {
					return answer
				}

				exp := regexp.MustCompile(`[,\s]+`)
				return exp.Split(domains, -1)
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

	return next(context.WithValue(
		ctx,
		"project",
		entities.Project{
			Name:    answers.Name,
			Domains: answers.Domains,
			Volumes: []string{
				fmt.Sprintf("%s__wp", projectSlug),
			},
		},
	))
}

func newProjectMarshalPipe(ctx context.Context, next pipelines.NextPipe) error {
	data, err := yaml.Marshal(ctx.Value("project"))
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return next(ctx)
}
