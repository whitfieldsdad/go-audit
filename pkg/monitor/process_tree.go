package monitor

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
