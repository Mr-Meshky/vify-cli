package cmd

import (
	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Stops the VPN/proxy and cleanly restores system DNS and routing",
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}
		return application.RunDisconnect()
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
