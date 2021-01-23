package cmd

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"regexp"
	"strings"
	"wpld/templates"
)

var newQuestions = []*survey.Question{
	{
		Name: "Name",
		Validate: survey.Required,
		Prompt: &survey.Input{
			Message: "What is the title of your site?",
		},
	},
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
			Name string
			Hostname string
		}{}

		if err := survey.Ask(newQuestions, &answers); err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil
			}
			return err
		}

		re := regexp.MustCompile(`\W+`)
		basePath := strings.TrimSuffix(
			strings.ToLower(
				string(
					re.ReplaceAll(
						[]byte(answers.Name),
						[]byte("-"),
					),
				),
			),
			"-",
		)

		wpldFilepath := "wpld.yaml"
		nginxTemplateFilepath := filepath.FromSlash("config/nginx/default.conf.template")

		compose := viper.New()

		compose.Set("name", answers.Name)
		compose.Set("hostname", answers.Hostname)

		compose.Set("services.cache.image", "memcached:latest")

		compose.Set("services.nginx.image", "nginx:latest")
		compose.Set("services.nginx.volumes", []string{
			nginxTemplateFilepath + ":/etc/nginx/templates",
		})

		compose.Set("services.wordpress.image", "wordpress:latest")

		config, err := yaml.Marshal(compose.AllSettings())
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()
		files := map[string][]byte {
			wpldFilepath: config,
			nginxTemplateFilepath: []byte(templates.NGINX_TEMPLATE),
		}

		for filename, data := range files {
			path := filepath.Join(basePath, filename)

			if mkdirErr := fs.MkdirAll(filepath.Dir(path), 0755); mkdirErr != nil {
				return mkdirErr
			}

			if writeErr := afero.WriteFile(fs, path, data, 0644); writeErr != nil {
				return writeErr
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
