package monitor

import (
	"github.com/charmbracelet/log"
	"github.com/hashicorp/golang-lru/v2/expirable"
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

type ProcessRef struct {
	Pid  int  `json:"pid"`
	Ppid *int `json:"ppid,omitempty"`
}

type Process struct {
	Name       string `json:"name"`
	Pid        int    `json:"pid"`
	Ppid       *int   `json:"ppid,omitempty"`
	Executable *File  `json:"executable,omitempty"`
}

type ProcessTreeInterface interface {
	AddProcess(process *Process) error
	RemoveProcess(pid int) error
	GetParentPid(pid int) (*int, error)
	GetChildPids(pid int) ([]int, error)
	GetAncestorPids(pid int) ([]int, error)
	GetDescendantPids(pid int) ([]int, error)
	GetSiblingPids(pid int) ([]int, error)
	GetParent(pid int) *Process
	GetChildren(pid int) ([]Process, error)
	GetAncestors(pid int) ([]Process, error)
	GetDescendants(pid int) ([]Process, error)
	GetSiblings(pid int) ([]Process, error)
}

type ProcessTree struct {
}

func NewProcessTree() ProcessTree {
	return ProcessTree{}
}

func (t *ProcessTree) AddProcess(process *Process) error {
	panic("implement me")
}

func (t *ProcessTree) RemoveProcess(pid int) error {
	panic("implement me")
}

func (t *ProcessTree) GetParentPid(pid int) (*int, error) {
	panic("implement me")
}

func (t *ProcessTree) GetChildPids(pid int) ([]int, error) {
	panic("implement me")
}

func (t *ProcessTree) GetAncestorPids(pid int) ([]int, error) {
	panic("implement me")
}

func (t *ProcessTree) GetDescendantPids(pid int) ([]int, error) {
	panic("implement me")
}

func (t *ProcessTree) GetSiblingPids(pid int) ([]int, error) {
	panic("implement me")
}

func (t *ProcessTree) GetParent(pid int) *Process {
	panic("implement me")
}

func (t *ProcessTree) GetChildren(pid int) ([]Process, error) {
	panic("implement me")
}

func (t *ProcessTree) GetAncestors(pid int) ([]Process, error) {
	panic("implement me")
}

func (t *ProcessTree) GetDescendants(pid int) ([]Process, error) {
	panic("implement me")
}

func (t *ProcessTree) GetSiblings(pid int) ([]Process, error) {
	panic("implement me")
}
