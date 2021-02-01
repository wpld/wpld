package cmd

import (
	"github.com/spf13/cobra"
	"wpld/cases"
	"wpld/global"
	"wpld/models"
	"wpld/utils"
)

var downCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Use:          "down",
	Short:        "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runDown,
}

func init() {
	rootCmd.AddCommand(downCmd)

	flags := downCmd.Flags()
	flags.Bool("rm", false, "Remove containers.")
	flags.Bool("all", false, "Down all containers.")
}

func runDown(cmd *cobra.Command, _ []string) error {
	config, err := cases.GetProjectConfig()
	if err != nil {
		return err
	}

	cli, err := cases.GetDockerClient()
	if err != nil {
		return err
	}

	factory := models.NewDockerFactory(cmd.Context(), cli)

	rm := false
	if rmFlag, rmErr := cmd.Flags().GetBool("rm"); rmErr == nil {
		rm = rmFlag
	}

	all := false
	if allFlag, allErr := cmd.Flags().GetBool("all"); allErr == nil {
		all = allFlag
	}

	if err = cases.StopArbitraryContainer(factory, config.Sub("nginx"), rm); err != nil {
		return err
	}

	if err = cases.StopArbitraryContainer(factory, config.Sub("wordpress"), rm); err != nil {
		return err
	}

	prefix := utils.Slugify(config.GetString("name"))
	services := config.Sub("services")
	for key := range services.AllSettings() {
		service := services.Sub(key)
		service.SetDefault("name", prefix+"_"+key)
		if err = cases.StopArbitraryContainer(factory, service, rm); err != nil {
			return err
		}
	}

	if !all {
		return nil
	}

	stopContainers := []func(models.DockerFactory, bool) error{
		global.StopMyAdmin,
		global.StopMySQL,
		global.StopDnsMasq,
		global.StopNginxProxy,
	}

	for _, stopContainer := range stopContainers {
		if stopErr := stopContainer(factory, rm); stopErr != nil {
			return stopErr
		}
	}

	return nil
}
