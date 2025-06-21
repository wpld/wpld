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

func ProjectPromptPipe() pipelines.Pipe {
	return func(ctx context.Context, next pipelines.NextPipe) error {
		var answers struct {
			Name    string
			Type    string
			Domains string
			PHP     string `survey:"php"`
			Cache   string
			// WordPressType  string `survey:"wordpress-type"`
			// WordPressUser  string `survey:"wordpress-user"`
			// WordPressPass  string `survey:"wordpress-pass"`
			// WordPressEmail string `survey:"wordpress-email"`
		}

		questions := []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "Project name:",
				},
				Validate: survey.Required,
			},
			{
				Name: "domains",
				Prompt: &survey.Input{
					Message: "Domain names:",
				},
				Validate: survey.Required,
			},
			{
				Name: "type",
				Prompt: &survey.Select{
					Message: "Project Type:",
					Default: specs.PROJECT_TYPE_PLUGIN,
					Options: []string{
						specs.PROJECT_TYPE_PLUGIN,
						specs.PROJECT_TYPE_THEME,
						specs.PROJECT_TYPE_WP_CONTENT,
					},
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
			// {
			// 	Name: "wordpress-type",
			// 	Prompt: &survey.Select{
			// 		Message: "Select a WordPress installation type:",
			// 		Default: specs.WORDPRESS_TYPE_SINGLE,
			// 		Options: []string{
			// 			specs.WORDPRESS_TYPE_SINGLE,
			// 			specs.WORDPRESS_TYPE_SUBDIR,
			// 			specs.WORDPRESS_TYPE_SUBDOMAINS,
			// 		},
			// 	},
			// },
			// {
			// 	Name:     "wordpress-user",
			// 	Validate: survey.Required,
			// 	Prompt: &survey.Input{
			// 		Message: "Admin Username",
			// 		Default: "admin",
			// 	},
			// },
			// {
			// 	Name:     "wordpress-pass",
			// 	Validate: survey.Required,
			// 	Prompt: &survey.Input{
			// 		Message: "Admin Password",
			// 		Default: "password",
			// 	},
			// },
			// {
			// 	Name:     "wordpress-email",
			// 	Validate: survey.Required,
			// 	Prompt: &survey.Input{
			// 		Message: "Admin Email",
			// 		Default: "admin@example.com",
			// 	},
			// },
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
			"wp":    specs.NewWordPressSpec(projectSlug, projectSlug, wpVolume, answers.PHP, answers.Type),
			"db":    specs.NewDatabaseSpec(projectSlug, dbVolume),
			"nginx": specs.NewNginxSpec(regexp.MustCompile(`[,\s]+`).Split(answers.Domains, -1)),
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

		// wp := entities.WordPress{
		// 	User:         answers.WordPressUser,
		// 	Password:     answers.WordPressPass,
		// 	Email:        answers.WordPressEmail,
		// 	Multisite:    false,
		// 	Subdirectory: false,
		// }
		//
		// if answers.WordPressType == specs.WORDPRESS_TYPE_SUBDIR {
		// 	wp.Subdirectory = true
		// } else if answers.WordPressType == specs.WORDPRESS_TYPE_SUBDOMAINS {
		// 	wp.Multisite = true
		// }

		ctx = context.WithValue(
			ctx,
			"project",
			entities.Project{
				ID:       projectSlug,
				Name:     answers.Name,
				Volumes:  volumes,
				Services: services,
				// WP:       wp,
			},
		)

		return next(ctx)
	}
}
