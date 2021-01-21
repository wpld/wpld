package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"strings"
	"wpld/global"
	"wpld/models"
)

var newQuestions = []*survey.Question{
	{
		Name: "Hostname",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "What is the primary hostname for your site? (Ex: docker.test)",
		},
	},
}

var newCmd = &cobra.Command{
	SilenceUsage: true,
	Args: cobra.NoArgs,
	Use: "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Hostname string
		}{}

		if err := survey.Ask(newQuestions, &answers); err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil
			}
			return err
		}

		compose := models.Compose{
			Version: "2.4",
			Services: map[string]models.Service{
				"cache": {
					Image: "memcached:latest",
					Networks: []string{
						global.NETWORK_NAME,
					},
				},
				"nginx": {
					Image: "nginx:latest",
					Expose: []int{
						80,
						443,
					},
					Volumes: []string{
						"./wordpress:/var/www/html:cached",
						"./config/nginx/default.conf.template:/etc/nginx/templates",
					},
					Networks: []string{
						global.NETWORK_NAME,
					},
					DependsOn: []string{
						"phpfpm",
					},
					Environment: map[string]string{
						"CERT_NAME": "",
						"HTTPS_METHOD": "noredirect",
						"VIRTUAL_HOST": fmt.Sprintf("%[1]s,*.%[1]s", answers.Hostname),
					},
				},
				"phpfpm": {
					Image: "wordpress",
					Networks: []string{
						global.NETWORK_NAME,
					},
					DependsOn: []string{
						"cache",
					},
					Environment: map[string]string{
						"ENABLE_XDEBUG": "true",
						"WORDPRESS_DB_HOST": global.MYSQL_CONTAINER_NAME,
						"WORDPRESS_DB_USER": "wordpress",
						"WORDPRESS_DB_PASSWORD": "password",
						"WORDPRESS_DB_NAME": strings.ReplaceAll(answers.Hostname, ".", "-"),
					},
				},
			},
			Networks: map[string]models.Network{
				global.NETWORK_NAME: {
					External: map[string]string{
						"name": global.NETWORK_NAME,
					},
				},
			},
		}

		fs := afero.NewOsFs()

		if err := compose.Save(fs); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
