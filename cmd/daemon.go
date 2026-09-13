package cmd

import (
	"fmt"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/Mr-Meshky/vify-cli/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonPort int

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run Vify background daemon providing IPC and local API for Desktop/Mobile GUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}

		server := daemon.NewServer(application)
		fmt.Printf("⚡ Vify Daemon API listening on 127.0.0.1:%d\n", daemonPort)
		return server.Start(daemonPort)
	},
}

func init() {
	daemonCmd.Flags().IntVarP(&daemonPort, "port", "P", 28190, "Port for daemon REST/IPC API")
	rootCmd.AddCommand(daemonCmd)
}
