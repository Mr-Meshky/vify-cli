package cmd

import (
	"context"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var cleanIPCount int

var cleanIPCmd = &cobra.Command{
	Use:   "clean-ip",
	Short: "Scans and finds optimal Cloudflare Clean IPs for your current ISP (MCI, Irancell, TCI, etc.)",
	Example: `  vify clean-ip
  vify clean-ip --count 15`,
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}
		return application.RunCleanIP(context.Background(), cleanIPCount)
	},
}

func init() {
	cleanIPCmd.Flags().IntVarP(&cleanIPCount, "count", "n", 10, "Number of top clean IPs to display")
	rootCmd.AddCommand(cleanIPCmd)
}
