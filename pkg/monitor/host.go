package monitor

import (
	"os"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/denisbrodbeck/machineid"
)

const (
	appId = "a86148b4-19e2-4533-a16f-3f3e96e92848"
)

type Host struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	OS       OS     `json:"os"`
}

type OS struct {
	Type string `json:"type"`
	Arch string `json:"arch"`
}

func GetOS() OS {
	return OS{
		Type: runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func GetHost() Host {
	id := GetHostId()
	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Failed to get hostname: %s", hostname)
	}
	return Host{
		Id:       id,
		Hostname: hostname,
		OS:       GetOS(),
	}
}

var _hostId string

func GetHostId() string {
	if _hostId == "" {
		id, err := machineid.ID()
		if err != nil {
			log.Fatalf("Failed to get host ID")
		}
		_hostId = NewUUID5(appId, []byte(id))
	}
	return _hostId
}
