package cmd

import (
	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		config := viper.New()
		config.SetFs(afero.NewOsFs())
		config.SetConfigName("wpld")
		config.SetConfigType("yaml")
		config.AddConfigPath(dir)
		if err = config.ReadInConfig(); err != nil {
			return err
		}

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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

		// TODO: stop containers using goroutines
		if err = global.StopMyAdmin(factory, rm); err != nil {
			return err
		}

		if err = global.StopMySQL(factory, rm); err != nil {
			return err
		}

		if err = global.StopDnsMasq(factory, rm); err != nil {
			return err
		}

		if err = global.StopNginxProxy(factory, rm); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	flags := downCmd.Flags()
	flags.Bool("rm", false, "Remove containers.")
	flags.Bool("all", false, "Down all containers.")
}
