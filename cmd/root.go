package cmd

import (
	"context"

	"github.com/charmbracelet/log"
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
		debug, _ := cmd.Flags().GetBool("debug")
		setLogLevel(debug)

		ctx := context.Background()
		monitor, err := monitor.NewAuditMonitor(nil)
		if err != nil {
			log.Fatalf("Failed to create process monitor: %v", err)
		}
		err = monitor.Run(ctx)
		if err != nil {
			log.Fatalf("Failed to run process monitor: %v", err)
		}
	},
}

func setLogLevel(debug bool) {
	var level log.Level
	if debug {
		level = log.DebugLevel
	} else {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug")
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
