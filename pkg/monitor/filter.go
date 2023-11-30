package monitor

import (
	"slices"

	"github.com/mitchellh/go-ps"
)

type Filter interface {
	AddByPid(pids []int) error
	AddByPath(paths []string) error
}

type ProcessFilter struct {
	Filter
	Pids  []int
	Paths []string
}

func NewProcessFilter() *ProcessFilter {
	return &ProcessFilter{}
}

func (f *ProcessFilter) AddByPid(pids []int) error {
	processes, err := ps.Processes()
	if err != nil {
		return err
	}
	for _, process := range processes {
		pid := process.Pid()
		if len(pids) > 0 && !slices.Contains(pids, pid) {
			continue
		}
		path := process.Executable()
		if len(f.Paths) > 0 && !slices.Contains(f.Paths, path) {
			continue
		}
		f.Pids = append(f.Pids, pid)
	}
	return nil
}

func (f *ProcessFilter) AddByPath(paths []string) error {
	f.Paths = append(f.Paths, paths...)
	processes, err := ps.Processes()
	if err != nil {
		return err
	}
	for _, process := range processes {
		path := process.Executable()
		if len(paths) > 0 && !slices.Contains(paths, path) {
			continue
		}
		f.Pids = append(f.Pids, process.Pid())
	}
	return nil
}
