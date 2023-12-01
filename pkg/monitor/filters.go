package monitor

import "github.com/whitfieldsdad/go-audit/pkg/util"

type ProcessFilterInterface interface {
	AddPids(pids ...int)
	RemovePids(pids ...int)
	AddPpids(ppids ...int)
	RemovePpids(ppids ...int)
	Merge(other ProcessFilter)
}

type ProcessFilter struct {
	Pids  []int `json:"pids,omitempty"`
	Ppids []int `json:"ppids,omitempty"`
}

func NewProcessFilter() *ProcessFilter {
	return &ProcessFilter{}
}

func (f *ProcessFilter) Merge(other ProcessFilter) {
	f.AddPids(other.Pids...)
	f.AddPpids(other.Ppids...)
}

func (f *ProcessFilter) AddPids(pids ...int) {
	f.Pids = util.AddIntsToSet(f.Pids, pids)
}

func (f *ProcessFilter) RemovePids(pids ...int) {
	f.Pids = util.RemoveIntsFromSet(f.Pids, pids)
}

func (f *ProcessFilter) AddPpids(ppids ...int) {
	f.Ppids = util.AddIntsToSet(f.Ppids, ppids)
}

func (f *ProcessFilter) RemovePpids(ppids ...int) {
	f.Ppids = util.RemoveIntsFromSet(f.Ppids, ppids)
}
