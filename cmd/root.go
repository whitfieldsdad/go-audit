package cmd

import (
	"context"
	"os"
	"os/signal"

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

		monitor, err := monitor.NewAuditMonitor()
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
