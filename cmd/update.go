package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mr-Meshky/vify-cli/internal/tui"
	"github.com/Mr-Meshky/vify-cli/internal/updater"
	"github.com/spf13/cobra"
)

var (
	updateCheckOnly bool
	updateForce     bool

	updateCmd = &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade"},
		Short:   "Check for and install the latest vify release",
		Long:    tui.TitleStyle.Render("⚡ VIFY CLI — Self Updater") + "\n\nChecks GitHub Releases for new versions and upgrades the CLI in-place.",
		Run: func(cmd *cobra.Command, args []string) {
			current := version
			fmt.Printf("%s Current version: %s\n", tui.TitleStyle.Render("⚡ Vify:"), tui.SuccessBadge.Render("v"+strings.TrimPrefix(current, "v")))

			if updateCheckOnly {
				fmt.Println("Checking GitHub for the latest release...")
				release, isNewer, err := updater.CheckUpdate(current, true)
				if err != nil {
					fmt.Printf("%s %v\n", tui.DangerBadge.Render("Error:"), err)
					os.Exit(1)
				}

				if isNewer {
					fmt.Printf("\n%s\n", updater.RenderUpdateNotification(current, release.TagName))
					if release.Body != "" {
						fmt.Printf("%s\n%s\n", tui.SubtitleStyle.Render("Changelog:"), strings.TrimSpace(release.Body))
					}
				} else {
					fmt.Printf("\n%s You are already on the latest version (%s)\n", tui.SuccessBadge.Render("✓ Up to date:"), release.TagName)
				}
				return
			}

			// Perform self-update
			release, err := updater.SelfUpdate(current, func(step string) {
				fmt.Printf("%s %s\n", tui.BadgeStyle.Render("→"), step)
			})

			if err != nil {
				fmt.Println()
				if strings.Contains(err.Error(), "permission denied") {
					fmt.Printf("%s %v\n", tui.DangerBadge.Render("Permission Denied:"), err)
					fmt.Println("Tip: Try running the update with administrator privileges (e.g. 'sudo vify update').")
				} else {
					fmt.Printf("%s %v\n", tui.DangerBadge.Render("Update Failed:"), err)
				}
				os.Exit(1)
			}

			if release != nil && updater.CompareVersions(release.TagName, current) <= 0 && !updateForce {
				fmt.Printf("\n%s You are already running the latest version (%s)\n", tui.SuccessBadge.Render("✓ Up to date:"), release.TagName)
				return
			}

			fmt.Println()
			fmt.Printf("%s Successfully updated to %s!\n", tui.SuccessBadge.Render("SUCCESS"), release.TagName)
			fmt.Printf("Run %s to see what's new.\n", tui.TitleStyle.Render("vify --version"))
		},
	}
)

func init() {
	updateCmd.Flags().BoolVarP(&updateCheckOnly, "check", "c", false, "Only check for new releases without installing")
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Force download and reinstall latest version even if up to date")
	rootCmd.AddCommand(updateCmd)
}
