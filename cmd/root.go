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

		// Build the process filter
		pids, _ := cmd.Flags().GetIntSlice("pids")
		ppids, _ := cmd.Flags().GetIntSlice("ppids")

		var filter *monitor.ProcessFilter
		if len(pids) > 0 || len(ppids) > 0 {
			filter = &monitor.ProcessFilter{
				Pids:  pids,
				Ppids: ppids,
			}
		}

		// Run the monitor
		ctx := context.Background()
		monitor, err := monitor.NewAuditMonitor(filter)
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
	runCmd.Flags().IntSliceP("pids", "p", []int{}, "PIDs")
	runCmd.Flags().IntSlice("ppids", []int{}, "PPIDs")
	runCmd.Flags().IntSlice("ancestor-pids", []int{}, "Ancestor PIDs")
}

func Execute() error {
	return rootCmd.Execute()
}
