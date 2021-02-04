package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
	"wpld/cases"
	"wpld/compose"
)

var buildCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Use:          "build",
	Short:        "A brief description of your command",
	RunE:         runBuild,
}

func init() {
	buildCmd.Long = heredoc.Doc(`
		A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.
	`)

	rootCmd.AddCommand(buildCmd)
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func runBuild(cmd *cobra.Command, args []string) error {
	configPath, err := filepath.Abs(".wpld/config.yaml")
	if err != nil {
		return err
	}

	fs := afero.NewOsFs()
	configData, err := afero.ReadFile(fs, configPath)
	if err != nil {
		return err
	}

	var config compose.Compose
	if err = yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	cli, err := cases.GetDockerClient()
	if err != nil {
		return err
	}

	for key, service := range config.Services {
		if service.Build.Name != "" {
			contextFolder := filepath.Dir(configPath)
			if service.Build.Context != "" {
				contextFolder = service.Build.Context
				if ! filepath.IsAbs(contextFolder) {
					contextFolder = filepath.Join(filepath.Dir(configPath), contextFolder)
				}
			}

			dockerFile := service.Build.Dockerfile
			if ! filepath.IsAbs(dockerFile) {
				dockerFile = filepath.Join(contextFolder, service.Build.Dockerfile)
			}

			dockerFile, dockerFileErr := filepath.Rel(contextFolder, dockerFile)
			if dockerFileErr != nil {
				return dockerFileErr
			}

			tag := service.Build.Name
			if tag == "" {
				tag = fmt.Sprintf("%s-%s", strings.ToLower(service.Name), strings.ToLower(key))
			}

			buildArgs := map[string]*string{}
			if len(service.Build.Args) > 0 {
				for argKey, argVal := range service.Build.Args {
					value := argVal
					buildArgs[ argKey ] = &value
				}
			}

			buildOptions := types.ImageBuildOptions{
				Dockerfile: dockerFile,
				Tags: []string{ tag },
				BuildArgs: buildArgs,
			}

			tar, tarErr := archive.TarWithOptions(contextFolder, &archive.TarOptions{})
			if tarErr != nil {
				return tarErr
			}

			buildRes, buildErr := cli.ImageBuild(cmd.Context(), tar, buildOptions)
			if buildErr != nil {
				return buildErr
			}

			scanner := bufio.NewScanner(buildRes.Body)
			for scanner.Scan() {
				lastLine := scanner.Text()

				errLine := &ErrorLine{}
				jsonErr := json.Unmarshal([]byte(lastLine), errLine)
				if jsonErr == nil && len(errLine.Error) > 0 {
					return errors.New(errLine.Error)
				}

				fmt.Println(lastLine)
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}
	}

	return nil
}
