package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"

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

		// Handle SIGINT and SIGTERM.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		go func() {
			<-sigCh
			log.Info("Shutting down...")
			os.Exit(0)
		}()

		f := &monitor.ProcessFilter{}
		ancestorPids, _ := cmd.Flags().GetInt32Slice("ancestor-pid")
		f.AncestorPIDs = ancestorPids

		monitor, err := monitor.NewAuditMonitor(f)
		if err != nil {
			log.Fatalf("Failed to create process monitor: %v", err)
		}

		// Run the monitor in a goroutine.
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := monitor.Run(ctx)
			if err != nil {
				log.Errorf("Monitor failed: %v", err)
			}
		}()

		// Continuously read events from the monitor.
		for {
			select {
			case event := <-monitor.Events:
				b, err := json.Marshal(event)
				if err != nil {
					log.Fatalf("Failed to marshal event: %v", err)
				}
				fmt.Println(string(b))
			}
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
	runCmd.PersistentFlags().Int32Slice("ancestor-pid", []int32{}, "Ancestor PIDs")

	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
