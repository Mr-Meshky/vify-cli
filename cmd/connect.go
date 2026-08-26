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
	connectCountry   string
	connectProtocol  string
	connectTun       bool
	connectSysProxy  bool
	connectFastPass  bool
	connectBatch     int
	connectCacheOnly bool
	connectRawURL    string
)

var connectCmd = &cobra.Command{
	Use:   "connect [CONFIG_URL]",
	Short: "Connect to VPN (Auto-select best server or pass manual vless/vmess/trojan/ss link)",
	Example: `  # Connect to VPN (Full TUN Mode - V2Box style)
  sudo vify connect

  # Connect to a manual config link directly
  sudo vify connect "trojan://password@44.246.163.102:443?security=tls&sni=example.com#MyServer"

  # Connect with filters
  sudo vify connect --country DE --protocol vless`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}

		rawLink := connectRawURL
		if len(args) > 0 && args[0] != "" {
			rawLink = args[0]
		}

		mode := "tun" // Default to Full TUN VPN like V2Box
		if connectSysProxy || (cmd.Flags().Changed("tun") && !connectTun) {
			mode = "system_proxy"
		}

		// TUN mode needs root to create a virtual network interface. Check
		// this upfront rather than letting the user sit through a full
		// subscription fetch + benchmark only to hit a cryptic sing-box
		// "operation not permitted" error at the very last step.
		if mode == "tun" && runtime.GOOS != "windows" && os.Geteuid() != 0 {
			return fmt.Errorf("TUN mode requires root privileges (creating a virtual network interface) — rerun with sudo, or use --system-proxy to connect without root")
		}

		return application.RunConnect(context.Background(), app.ConnectOptions{
			Country:      connectCountry,
			Protocol:     connectProtocol,
			Mode:         mode,
			FastPass:     connectFastPass,
			BatchSize:    connectBatch,
			UseCacheOnly: connectCacheOnly,
			ManualURI:    rawLink,
		})
	},
}

func init() {
	connectCmd.Flags().StringVarP(&connectCountry, "country", "c", "", "Filter candidate servers by country code (e.g., DE, NL, US, TR)")
	connectCmd.Flags().StringVarP(&connectProtocol, "protocol", "p", "", "Filter candidate servers by protocol (vless, vmess, trojan, ss)")
	connectCmd.Flags().StringVarP(&connectRawURL, "url", "u", "", "Connect directly to a specific vless/vmess/trojan/ss link")
	connectCmd.Flags().BoolVar(&connectTun, "tun", true, "Use TUN mode (Full virtual network interface / device VPN like V2Box)")
	connectCmd.Flags().BoolVar(&connectSysProxy, "system-proxy", false, "Use System Proxy mode instead of TUN")
	connectCmd.Flags().BoolVar(&connectFastPass, "fast-pass", true, "Connect immediately upon finding the first healthy sub-800ms config")
	connectCmd.Flags().IntVarP(&connectBatch, "batch", "b", 60, "Maximum number of candidate servers to test in parallel")
	connectCmd.Flags().BoolVar(&connectCacheOnly, "cache-only", false, "Use locally cached servers without fetching latest subscriptions")

	rootCmd.AddCommand(connectCmd)
}
