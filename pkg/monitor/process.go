package monitor

import (
	"fmt"
	"time"

	ps "github.com/shirou/gopsutil/v3/process"
)

type ProcessOptions struct {
	IncludeHashes bool `json:"include_hashes"`
}

func GetDefaultProcessOptions() *ProcessOptions {
	return &ProcessOptions{
		IncludeHashes: true,
	}
}

type Process struct {
	PID         int32      `json:"pid"`
	PPID        int32      `json:"ppid"`
	Name        string     `json:"name,omitempty"`
	Argv        []string   `json:"argv,omitempty"`
	Argc        int        `json:"argc,omitempty"`
	CommandLine string     `json:"command_line,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
	ExitCode    *int       `json:"exit_code,omitempty"`
	Executable  *File      `json:"executable,omitempty"`
}

func (p Process) Hash() uint64 {
	return calculateProcessId(p.PID, p.PPID)
}

func ListProcesses(opts *ProcessOptions) ([]Process, error) {
	if opts == nil {
		opts = GetDefaultProcessOptions()
	}
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	var results []Process
	for _, p := range processes {
		process := parseProcess(p)
		if opts.IncludeHashes && process.Executable != nil {
			hashes, err := GetFileHashes(process.Executable.Path)
			if err != nil {
				return nil, err
			}
			process.Executable.Hashes = hashes
		}
		results = append(results, process)
	}
	return results, nil
}

func GetProcess(pid int32, opts *ProcessOptions) (*Process, error) {
	if opts == nil {
		opts = GetDefaultProcessOptions()
	}
	p, err := ps.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	process := parseProcess(p)
	if opts.IncludeHashes && process.Executable != nil {
		hashes, err := GetFileHashes(process.Executable.Path)
		if err != nil {
			return nil, err
		}
		process.Executable.Hashes = hashes
	}
	return &process, nil
}

func parseProcess(p *ps.Process) Process {
	pid := p.Pid
	ppid, _ := p.Ppid()
	name, _ := p.Name()

	var executable *File
	executablePath, _ := p.Exe()
	if executablePath != "" {
		executable, _ = GetFile(executablePath)
	}

	argv, _ := p.CmdlineSlice()
	argc := len(argv)
	commandLine, _ := p.Cmdline()

	var createTime *time.Time
	createTimeMs, err := p.CreateTime()
	if err == nil {
		t := time.UnixMilli(createTimeMs)
		createTime = &t
	}
	process := Process{
		PID:         pid,
		PPID:        ppid,
		Name:        name,
		Argv:        argv,
		Argc:        argc,
		CommandLine: commandLine,
		CreateTime:  createTime,
		Executable:  executable,
	}
	return process
}

type ProcessIdentity struct {
	PID  int32 `json:"pid"`
	PPID int32 `json:"ppid"`
}

func (p ProcessIdentity) Hash() uint64 {
	return calculateProcessId(p.PID, p.PPID)
}

func calculateProcessId(pid, ppid int32) uint64 {
	k := []byte(fmt.Sprintf("%d,%d", pid, ppid))
	return GetXXH3(k)
}
