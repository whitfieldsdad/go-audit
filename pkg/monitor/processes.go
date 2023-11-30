package monitor

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
)

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
