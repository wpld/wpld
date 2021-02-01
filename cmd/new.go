package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path/filepath"
	"wpld/compose"
	"wpld/global"
	"wpld/templates"
	"wpld/utils"
)

var newQuestions = []*survey.Question{
	{
		Name:     "Name",
		Validate: survey.Required,
		Prompt:   &survey.Input{
			Message: "What is the title of your site?",
		},
	},
	{
		Name:     "Hostname",
		Validate: survey.Required,
		Prompt:   &survey.Input{
			Message: "What is the primary hostname for your site? (Ex: docker.test)",
		},
	},
}

var newCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Use:          "new",
	Short:        "A brief description of your command",
	RunE:         runNew,
}

func init() {
	newCmd.Long = heredoc.Doc(`
		A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.
	`)

	rootCmd.AddCommand(newCmd)
}

func runNew(_ *cobra.Command, _ []string) error {
	var config compose.Compose
	if err := survey.Ask(newQuestions, &config); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil
		}
		return err
	}

	slug := utils.Slugify(config.Name)

	wpldFilepath := filepath.FromSlash(".wpld/config.yaml")
	phpDockerfileFilepath := filepath.FromSlash(".wpld/php/Dockerfile")
	nginxTemplateFilepath := filepath.FromSlash(".wpld/nginx/default.conf.template")

	config.Services = make(map[string]compose.Service)

	config.Services["cache"] = compose.Service{
		Image: "memcached:latest",
		Name: fmt.Sprintf("%s_cache", slug),
	}

	config.Services["nginx"] = compose.Service{
		Image: "nginx:latest",
		Name: fmt.Sprintf("%s_nginx", slug),
		Volumes: []string{
			"nginx/default.conf.template:/etc/nginx/templates/default.conf.template",
			"wordpress:/var/www/html",
		},
		Env: map[string]string{
			"HTTPS_METHOD": "noredirect",
			"VIRTUAL_HOST": config.Hostname,
			"CERT_NAME":    slug,
			"PHPFPM_HOST":  fmt.Sprintf("%s_wordpress", slug),
		},
	}

	config.Services["wordpress"] = compose.Service{
		Name: fmt.Sprintf("%s_wordpress", slug),
		Build: compose.Build{
			Dockerfile: "Dockerfile",
			Context: "php",
			Args: map[string]string{
				"PHP_IMAGE": "8.0-fpm-alpine",
				"CALLING_USER": "",
				"CALLING_UID": "",
			},
		},
		Volumes: []string{
			"wordpress:/var/www/html",
		},
		Env: map[string]string{
			"WORDPRESS_DB_HOST":     global.MYSQL_CONTAINER_NAME,
			"WORDPRESS_DB_USER":     "root",
			"WORDPRESS_DB_PASSWORD": "password",
			"WORDPRESS_DB_NAME":     slug,
		},
	}

	fs := afero.NewOsFs()
	files := map[string][]byte{
		wpldFilepath:          config.Serialize(),
		phpDockerfileFilepath: []byte(templates.PHP_DOCKERFILE),
		nginxTemplateFilepath: []byte(templates.NGINX_TEMPLATE),
	}

	for filename, data := range files {
		path := filepath.Join(slug, filename)
		if mkdirErr := fs.MkdirAll(filepath.Dir(path), 0755); mkdirErr != nil {
			return mkdirErr
		}

		if writeErr := afero.WriteFile(fs, path, data, 0644); writeErr != nil {
			return writeErr
		}
	}

	return nil
}
