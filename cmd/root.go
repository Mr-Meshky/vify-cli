package cmd

import (
	"fmt"
	"os"

	"github.com/Mr-Meshky/vify-cli/internal/tui"
	"github.com/Mr-Meshky/vify-cli/internal/updater"
	"github.com/spf13/cobra"
)

var (
	version = "1.2.5"
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "vify",
		Short: "⚡ Vify CLI — Blazing-fast, lightweight, cross-platform terminal VPN client",
		Long:  tui.TitleStyle.Render("⚡ VIFY CLI"),
	}
)

// Execute runs the root CLI command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(tui.DangerBadge.Render("Error:") + " " + err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vify/config.yaml)")
	rootCmd.Version = version

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "update" || cmd.Name() == "help" {
			return
		}
		// Non-blocking check using local cache (only hits network if cache expired)
		if release, isNewer, err := updater.CheckUpdate(version, false); err == nil && isNewer && release != nil {
			fmt.Println(updater.RenderUpdateNotification(version, release.TagName))
		}
	}
}

