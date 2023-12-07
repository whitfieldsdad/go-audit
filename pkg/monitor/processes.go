package monitor

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

type Process struct {
	Name        string     `json:"name,omitempty"`
	Pid         int        `json:"pid"`
	Ppid        *int       `json:"ppid,omitempty"`
	Executable  *File      `json:"executable,omitempty"`
	Argv        []string   `json:"argv,omitempty"`
	Argc        int        `json:"argc,omitempty"`
	CommandLine string     `json:"command_line,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	ExitCode    *int       `json:"exit_code,omitempty"`
}

func GetProcess(pid int) (*Process, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	ppid32, err := p.Ppid()
	if err != nil {
		return nil, err
	}
	ppid := int(ppid32)

	var file *File
	exe, _ := p.Exe()
	if exe != "" {
		file, _ = GetFile(exe)
	}
	name, _ := p.Name()
	argv, _ := p.CmdlineSlice()
	argc := len(argv)
	cmdline := strings.Join(argv, " ")

	var (
		startTime *time.Time
	)
	startTimeMs, err := p.CreateTime()
	if err == nil {
		startTime = util.TimeFromMs(startTimeMs)
	}
	return &Process{
		Name:        name,
		Pid:         pid,
		Ppid:        &ppid,
		Executable:  file,
		Argv:        argv,
		Argc:        argc,
		CommandLine: cmdline,
		CreateTime:  startTime,
	}, nil
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

func GetProcessExecutable(pid int) (*File, error) {
	path, err := GetProcessExecutablePath(pid)
	if err != nil {
		return nil, nil
	}
	return GetFile(path)
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

func GetCommandLineArgs(pid int) ([]string, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process")
	}
	return p.CmdlineSlice()
}
