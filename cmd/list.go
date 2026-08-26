package cmd

import (
	"context"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	listCountry  string
	listProtocol string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Interactive TUI list to manually select and connect to a server",
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}
		return application.RunList(context.Background(), listCountry, listProtocol)
	},
}

func init() {
	listCmd.Flags().StringVarP(&listCountry, "country", "c", "", "Filter servers by country code (e.g. DE, NL, US)")
	listCmd.Flags().StringVarP(&listProtocol, "protocol", "p", "", "Filter servers by protocol (vless, vmess, trojan, ss)")
	rootCmd.AddCommand(listCmd)
}
