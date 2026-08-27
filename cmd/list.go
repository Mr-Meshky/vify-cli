package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	listCountry  string
	listProtocol string
	listSysProxy bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Interactive TUI list to manually pick which server and protocol to connect to",
	Example: `  vify list
  vify list --protocol vless
  vify list --country DE --system-proxy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}

		// Same root check as `connect`: picking a server here still opens a
		// full TUN VPN session by default, which needs a virtual network
		// interface. Fail fast with a clear message instead of letting the
		// user sit through selection only to hit a permissions error.
		if !listSysProxy && runtime.GOOS != "windows" && os.Geteuid() != 0 {
			return fmt.Errorf("TUN mode requires root privileges (creating a virtual network interface) — rerun with sudo, or use --system-proxy to connect without root")
		}

		return application.RunList(context.Background(), listCountry, listProtocol, listSysProxy)
	},
}

func init() {
	listCmd.Flags().StringVarP(&listCountry, "country", "c", "", "Filter servers by country code (e.g. DE, NL, US)")
	listCmd.Flags().StringVarP(&listProtocol, "protocol", "p", "", "Filter servers by protocol (vless, vmess, trojan, shadowsocks/ss)")
	listCmd.Flags().BoolVar(&listSysProxy, "system-proxy", false, "Use System Proxy mode instead of TUN (no root required)")
	rootCmd.AddCommand(listCmd)
}
