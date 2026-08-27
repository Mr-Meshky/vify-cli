package cmd

import (
	"fmt"
	"os"

	"github.com/Mr-Meshky/vify-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
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
}
