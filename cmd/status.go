package cmd

import (
	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays current connection state, active IP, country, and duration",
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}
		return application.RunStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
