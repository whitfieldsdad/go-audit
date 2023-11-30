module github.com/whitfieldsdad/go-audit

go 1.21.4

require (
	github.com/charmbracelet/log v0.3.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

require (
	github.com/0xrawsec/golang-etw v1.6.2
	github.com/0xrawsec/golang-utils v1.3.1 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/lipgloss v0.9.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/google/uuid v1.4.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.23.10
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
)

// Replace github.com/Velocidex/etw with github.com/whitfieldsdad/etw
replace github.com/Velocidex/etw => ../etw
