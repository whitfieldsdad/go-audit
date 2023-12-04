package monitor

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

type ProcessFilterInterface interface {
	AddPids(pids ...int)
	RemovePids(pids ...int)
	AddPpids(ppids ...int)
	RemovePpids(ppids ...int)
	Merge(other ProcessFilter)
}

type ProcessFilter struct {
	Pids           []int `json:"pids,omitempty"`
	Ppids          []int `json:"ppids,omitempty"`
	AncestorPids   []int `json:"ancestor_pids,omitempty"`
	DescendantPids []int `json:"descendant_pids,omitempty"`
}

func (f *ProcessFilter) Matches(p ProcessRef, t *ProcessTree) (bool, error) {
	if len(f.Pids) > 0 && !slices.Contains(f.Pids, p.Pid) {
		return false, nil
	}
	if len(f.Ppids) > 0 {
		if p.Ppid == nil || !slices.Contains(f.Ppids, *p.Ppid) {
			return false, nil
		}
	}
	if len(f.AncestorPids) > 0 || len(f.DescendantPids) > 0 {
		if t == nil {
			return false, errors.New("a process tree is required when filtering by ancestor or descendant PIDs")
		}
		if len(f.AncestorPids) > 0 {
			panic("not implemented")
		}
		if len(f.DescendantPids) > 0 {
			panic("not implemented")
		}
	}
	return true, nil
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

func (f *ProcessFilter) AddAncestorPids(pids ...int) {
	f.AncestorPids = util.AddIntsToSet(f.AncestorPids, pids)
}

func (f *ProcessFilter) RemoveAncestorPids(pids ...int) {
	f.AncestorPids = util.RemoveIntsFromSet(f.AncestorPids, pids)
}

func (f *ProcessFilter) AddDescendantPids(pids ...int) {
	f.DescendantPids = util.AddIntsToSet(f.DescendantPids, pids)
}

func (f *ProcessFilter) RemoveDescendantPids(pids ...int) {
	f.DescendantPids = util.RemoveIntsFromSet(f.DescendantPids, pids)
}
