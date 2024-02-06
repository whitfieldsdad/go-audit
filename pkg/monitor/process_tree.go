package monitor

type ProcessTree struct {
	pidToPpid map[int32]int32
}

func NewProcessTree() *ProcessTree {
	return &ProcessTree{
		pidToPpid: make(map[int32]int32),
	}
}

func NewProcessTreeFromProcessIdentities(ids []ProcessIdentity) *ProcessTree {
	t := NewProcessTree()
	for _, id := range ids {
		t.pidToPpid[id.PID] = id.PPID
	}
	return t
}

func GetProcessTree() (*ProcessTree, error) {
	ids, err := listProcessIdentities()
	if err != nil {
		return nil, err
	}
	return NewProcessTreeFromProcessIdentities(ids), nil
}

func (t *ProcessTree) AddProcess(ppid, pid int32) error {
	t.pidToPpid[pid] = ppid
	return nil
}

func (t *ProcessTree) RemoveProcesses(pids ...int32) {
	for _, pid := range pids {
		delete(t.pidToPpid, pid)
	}
}

func (t ProcessTree) GetAncestorPids(pid int32) []int32 {
	ancestors := []int32{}
	for {
		ppid, ok := t.pidToPpid[pid]
		if !ok {
			break
		}
		ancestors = append(ancestors, ppid)
		pid = ppid
	}
	return ancestors
}

func (t ProcessTree) GetDescendantPids(pid int32) []int32 {
	var descendants []int32
	children := t.GetChildPids(pid)
	for _, child := range children {
		descendants = append(descendants, child)
		descendants = append(descendants, t.GetDescendantPids(child)...)
	}
	return descendants
}

func (t ProcessTree) GetParentPid(pid int32) (int32, bool) {
	ppid, ok := t.pidToPpid[pid]
	return ppid, ok
}

func (t ProcessTree) GetChildPids(pid int32) []int32 {
	children := []int32{}
	for p, parent := range t.pidToPpid {
		if pid == parent {
			children = append(children, p)
		}
	}
	return children
}

func (t ProcessTree) GetSiblingPids(pid int32) []int32 {
	ppid, ok := t.GetParentPid(pid)
	if !ok {
		return nil
	}
	var siblings []int32
	for _, child := range t.GetChildPids(ppid) {
		if child != pid {
			siblings = append(siblings, child)
		}
	}
	return siblings
}

func (t ProcessTree) IsDescendantOfAny(pid int32, pids []int32) bool {
	for _, p := range pids {
		if t.IsDescendantOf(pid, p) {
			return true
		}
	}
	return false
}

func (t ProcessTree) IsDescendantOf(descendantPid, pid int32) bool {
	descendants := t.GetDescendantPids(pid)
	for _, descendant := range descendants {
		if descendant == descendantPid {
			return true
		}
	}
	return false
}
