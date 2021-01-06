package cmd

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"wpld/global"
)

var upCmd = &cobra.Command{
	SilenceUsage: true,
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Debugf("Running {%s} command...", cmd.Use)

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
		if err = global.RunMySQL(ctx, cli, pull); err != nil {
			return err
		}

		// TODO: wait until we can connect to the MySQL server before starting phpMyAdmin?
		if err = global.RunMyAdmin(ctx, cli, pull); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	f := upCmd.Flags()
	f.Bool("pull", false, "Pull images before starting containers.")
}