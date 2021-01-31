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

		if err = global.VerifyNetwork(factory); err != nil {
			return err
		}

		pull := false
		if pullFlag, err := cmd.Flags().GetBool("pull"); err == nil {
			pull = pullFlag
		}

		// TODO: start containers using goroutines
		if err = global.RunNginxProxy(factory, pull); err != nil {
			return err
		}

		if err = global.RunDnsMasq(factory, pull); err != nil {
			return err
		}

		if err = global.RunMySQL(factory, pull); err != nil {
			return err
		}

		for i := 0; i < 12; i++ {
			if db, err := sql.Open("mysql", "root:password@/information_schema"); err != nil {
				logrus.Error(err)
			} else {
				if pingErr := db.Ping(); pingErr == nil {
					_ = db.Close()
					break
				}
			}

			time.Sleep(5 * time.Second)
		}

		// TODO: wait until we can connect to the MySQL server before starting phpMyAdmin?
		if err = global.RunMyAdmin(factory, pull); err != nil {
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
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	f := upCmd.Flags()
	f.Bool("pull", false, "Pull images before starting containers.")
}
