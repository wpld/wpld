package cmd

import (
	"github.com/spf13/cobra"
	"wpld/cases"
	"wpld/global"
	"wpld/models"
	"wpld/utils"
)

var upCmd = &cobra.Command{
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	Use:          "up",
	Short:        "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runUp,
}

func init() {
	rootCmd.AddCommand(upCmd)

	f := upCmd.Flags()
	f.Bool("pull", false, "Pull images before starting containers.")
}

func runUp(cmd *cobra.Command, _ []string) error {
	config, err := cases.GetProjectConfig()
	if err != nil {
		return err
	}

	cli, err := cases.GetDockerClient()
	if err != nil {
		return err
	}

	factory := models.NewDockerFactory(cmd.Context(), cli)

	if err = global.VerifyNetwork(factory); err != nil {
		return err
	}

	pull := false
	if pullFlag, err := cmd.Flags().GetBool("pull"); err == nil {
		pull = pullFlag
	}

	runContainers := []func(models.DockerFactory, bool) error{
		global.RunNginxProxy,
		global.RunDnsMasq,
		global.RunMySQL,
		global.RunMyAdmin,
	}

	for _, runContainer := range runContainers {
		if runErr := runContainer(factory, pull); runErr != nil {
			return runErr
		}
	}

	if err = global.WaitForMySQL(); err != nil {
		return err
	}

	if err = cases.StartArbitraryContainer(factory, config.Sub("wordpress"), pull); err != nil {
		return err
	}

	if err = cases.StartArbitraryContainer(factory, config.Sub("nginx"), pull); err != nil {
		return err
	}

	prefix := utils.Slugify(config.GetString("name"))
	services := config.Sub("services")
	for key := range services.AllSettings() {
		service := services.Sub(key)
		service.SetDefault("name", prefix+"_"+key)
		if err = cases.StartArbitraryContainer(factory, service, pull); err != nil {
			return err
		}
	}

	return nil
}
