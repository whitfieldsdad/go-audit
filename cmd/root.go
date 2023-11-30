package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/whitfieldsdad/go-audit/pkg/monitor"
)

var rootCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit monitor",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the audit monitor",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		monitor, err := monitor.NewWindowsProcessMonitor()
		if err != nil {
			log.Fatalf("Failed to create process monitor: %v", err)
		}
		err = monitor.Run(ctx)
		if err != nil {
			log.Fatalf("Failed to run process monitor: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
