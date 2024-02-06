package monitor

type ProcessFilter struct {
	PIDs         []int32 `json:"pids"`
	AncestorPIDs []int32 `json:"ancestor_pids"`
}

func (f ProcessFilter) IsEmpty() bool {
	return len(f.PIDs) == 0 && len(f.AncestorPIDs) == 0
}
