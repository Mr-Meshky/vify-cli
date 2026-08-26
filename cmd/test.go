package cmd

import (
	"context"

	"github.com/Mr-Meshky/vify-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	testBatchSize int
	testCountry   string
	testProtocol  string
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Benchmarks current configs with real HTTP/204 get and displays a latency leaderboard",
	Example: `  vify test
  vify test --batch 40 --country DE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := app.New()
		if err != nil {
			return err
		}
		return application.RunTest(context.Background(), testBatchSize, testCountry, testProtocol)
	},
}

func init() {
	testCmd.Flags().IntVarP(&testBatchSize, "batch", "b", 40, "Maximum number of servers to benchmark")
	testCmd.Flags().StringVarP(&testCountry, "country", "c", "", "Filter candidate servers by country code (e.g. DE, NL)")
	testCmd.Flags().StringVarP(&testProtocol, "protocol", "p", "", "Filter candidate servers by protocol (vless, vmess, trojan, ss)")
	rootCmd.AddCommand(testCmd)
}
