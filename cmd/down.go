package cmd

import (
	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"wpld/cases"
	"wpld/global"
	"wpld/utils"
)

var downCmd = &cobra.Command{
	SilenceUsage: true,
	Args: cobra.NoArgs,
	Use: "down",
	Short: "A brief description of your command",
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

		ctx := cmd.Context()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}

		rm := false
		if rmFlag, rmErr := cmd.Flags().GetBool("rm"); rmErr == nil {
			rm = rmFlag
		}

		prefix := utils.Slugify(config.GetString("name"))
		services := config.Sub("services")
		for key, _ := range services.AllSettings() {
			service := services.Sub(key)
			service.SetDefault("name", "wpld_" + prefix + "_" + key)

			ctrn := cases.CreateArbitraryContainer(service)
			if err = ctrn.Stop(ctx, cli); err != nil {
				return err
			}

			if rm {
				if rmErr := ctrn.Remove(ctx, cli); rmErr != nil {
					return rmErr
				}
			}
		}

		// TODO: stop containers using goroutines
		if err = global.StopMyAdmin(ctx, cli, rm); err != nil {
			return err
		}

		if err = global.StopMySQL(ctx, cli, rm); err != nil {
			return err
		}

		if err = global.StopDnsMasq(ctx, cli, rm); err != nil {
			return err
		}

		if err = global.StopNginxProxy(ctx, cli, rm); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	flags := downCmd.Flags()
	flags.Bool("rm", false, "Remove containers.")
}
