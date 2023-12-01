package monitor

import (
	"github.com/charmbracelet/log"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
)

type processEntry struct {
	Pid  int
	PPid *int
}

func GetPPID(pid int) (*int, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process")
	}
	ppid32, err := p.Ppid()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ppid")
	}
	ppid := int(ppid32)
	return &ppid, nil
}

func GetProcessExecutablePath(pid int) (string, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}
	exe, err := p.Exe()
	if err != nil {
		return "", nil
	}
	return exe, nil
}

func GetProcessExecutable(pid int) (*File, error) {
	path, err := GetProcessExecutablePath(pid)
	if err != nil {
		return nil, nil
	}
	return GetFile(path)
}

func GetPPIDMap() (*expirable.LRU[int, int], error) {
	log.Infof("Initializing PPID map...")
	m := expirable.NewLRU[int, int](PidMapMaxSize, nil, PidMapTTL)
	processes, err := process.Processes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get processes")
	}
	for _, p := range processes {
		pid := int(p.Pid)
		ppid, err := p.Ppid()
		if err != nil {
			continue
		}
		m.Add(pid, int(ppid))
	}
	log.Infof("Loaded %d processes into PPID map", m.Len())
	return m, nil
}
