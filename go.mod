module github.com/whitfieldsdad/go-audit

go 1.21.4

require (
	github.com/Velocidex/etw v0.0.0-20231115144702-0b885b292f0f
	github.com/charmbracelet/log v0.3.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.8.0
	golang.org/x/sys v0.13.0
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/lipgloss v0.9.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
)

// Replace github.com/Velocidex/etw with github.com/whitfieldsdad/etw
replace github.com/Velocidex/etw => ../etw
