package terminal

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"wpld/internal/docker"
	"wpld/internal/pipelines"
	"wpld/internal/tasks"
)

var startCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "start",
	Short:         "start short desc",
	Aliases: []string{
		"up",
	},
	RunE: func(c *cobra.Command, args []string) error {
		pull, err := c.Flags().GetBool("pull")
		if err != nil {
			return err
		}

		api, err := docker.NewDocker()
		if err != nil {
			return err
		}

		fs := afero.NewOsFs()

		pipeline := pipelines.NewPipeline(
			tasks.ProjectUnmarshalPipe(fs),
			tasks.ContainersStartPipe(api, pull),
			tasks.DNSReloadPipe(api, fs),
		)

		return pipeline.Run(c.Context())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("pull", "p", false, "force pulling images before starting containers")
}
