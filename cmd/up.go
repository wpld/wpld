package cmd

import (
	"database/sql"
	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
	"wpld/cases"
	"wpld/global"
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

		if _, err = global.VerifyNetwork(ctx, cli); err != nil {
			return err
		}

		pull := false
		if pullFlag, err := cmd.Flags().GetBool("pull"); err == nil {
			pull = pullFlag
		}

		// TODO: start containers using goroutines
		if err = global.RunNginxProxy(ctx, cli, pull); err != nil {
			return err
		}

		if err = global.RunDnsMasq(ctx, cli, pull); err != nil {
			return err
		}

		if err = global.RunMySQL(ctx, cli, pull); err != nil {
			return err
		}

		for i := 0; i < 12; i++ {
			if db, err := sql.Open("mysql", "root:password@/information_schema"); err != nil {
				logrus.Error(err)
			} else {
				if pingErr := db.Ping(); pingErr == nil {
					db.Close()
					break
				}
			}

			time.Sleep(5 * time.Second)
		}

		// TODO: wait until we can connect to the MySQL server before starting phpMyAdmin?
		if err = global.RunMyAdmin(ctx, cli, pull); err != nil {
			return err
		}

		wp := cases.CreateArbitraryContainer(config.Sub("wordpress"))
		if err = wp.Start(ctx, cli); err != nil {
			return err
		}

		nginx := cases.CreateArbitraryContainer(config.Sub("nginx"))
		if err = nginx.Start(ctx, cli); err != nil {
			return err
		}

		prefix := utils.Slugify(config.GetString("name"))
		services := config.Sub("services")
		for key, _ := range services.AllSettings() {
			service := services.Sub(key)
			service.SetDefault("name", prefix + "_" + key)

			ctrn := cases.CreateArbitraryContainer(service)
			if err = ctrn.Start(ctx, cli); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	f := upCmd.Flags()
	f.Bool("pull", false, "Pull images before starting containers.")
}
