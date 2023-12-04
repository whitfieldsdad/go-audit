package monitor

import (
	"os"
	"runtime"
)

type Host struct {
	Hostname string `json:"hostname,omitempty"`
	OS       OS     `json:"os,omitempty"`
}

type OS struct {
	Type string `json:"type,omitempty"`
	Arch string `json:"arch,omitempty"`
}

func GetHost() Host {
	hostname, _ := os.Hostname()
	return Host{
		Hostname: hostname,
		OS:       GetOS(),
	}
}

func GetOS() OS {
	return OS{
		Type: runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
