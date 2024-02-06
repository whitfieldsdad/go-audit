package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestGetProcessTree(t *testing.T) {
	_, err := GetProcessTree()
	assert.Nil(t, err, "Failed to get process tree")
}

func TestGetAncestorPidsOneBranch(t *testing.T) {

	// 1 -> 2 -> 3 -> 4
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)

	expected := []int{1, 2}
	result := tree.GetAncestorPids(3)
	slices.Sort(result)
	assert.Equal(t, expected, result, "Failed to identify ancestors")
}

func TestGetAncestorPidsMultipleBranches(t *testing.T) {

	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	expected := []int{1, 5, 8}
	result := tree.GetAncestorPids(9)
	slices.Sort(result)
	assert.Equal(t, expected, result, "Failed to identify ancestors")
}

func TestGetDescendantPids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	expected := []int{6, 7, 8, 9}
	result := tree.GetDescendantPids(5)
	slices.Sort(result)
	assert.Equal(t, expected, result, "GetDescendantPids() should return the correct descendant pids")
}

func TestGetSiblingPids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	expected := []int{8}
	result := tree.GetSiblingPids(6)
	slices.Sort(result)
	assert.Equal(t, expected, result, "GetSiblingPids() should return the correct sibling pids")
}

func TestGetParentPid(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	expected := 2
	result, ok := tree.GetParentPid(3)
	assert.True(t, ok, "GetParentPid() should return true if the pid exists")
	assert.Equal(t, expected, result, "GetParentPid() should return the correct parent pid")
}

func TestGetChildPids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	//      2 -> 5 -> 6
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(2, 5)
	tree.AddProcess(5, 6)

	expected := []int{3, 5}
	result := tree.GetChildPids(2)
	slices.Sort(result)
	assert.Equal(t, expected, result, "GetChildPids() should return the correct child pids")
}

func TestIsParent(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 2
	ppid := 1
	assert.True(t, tree.IsParent(pid, ppid))
}

func TestIsParentWithSamePids(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 1
	ppid := 1
	assert.False(t, tree.IsParent(pid, ppid))
}

func TestIsParentWithIncorrectParent(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 1
	ppid := 2
	assert.False(t, tree.IsParent(pid, ppid))
}

func TestIsParentWithInvalidParent(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 1
	ppid := 0
	assert.False(t, tree.IsParent(pid, ppid))
}

func TestIsChild(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 3
	ppid := 2
	assert.True(t, tree.IsChild(pid, ppid))
}

func TestIsChildWithSamePids(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 1
	ppid := 1
	assert.False(t, tree.IsChild(pid, ppid))
}

func TestIsChildWithIncorrectParent(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 2
	ppid := 3
	assert.False(t, tree.IsChild(pid, ppid))
}

func TestIsChildWithInvalidParent(t *testing.T) {
	// 1 -> 2 -> 3
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)

	pid := 1
	ppid := 0
	assert.False(t, tree.IsChild(pid, ppid))
}

func TestIsSibling(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 3
	siblingPid := 5
	assert.True(t, tree.IsSibling(pid, siblingPid))
}

func TestIsSiblingWithInvalidHierarchy(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 3
	siblingPid := 2
	assert.False(t, tree.IsSibling(pid, siblingPid))
}

func TestIsSiblingWithInvalidPid(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 0
	siblingPid := 5
	assert.False(t, tree.IsSibling(pid, siblingPid))
}

func TestIsSiblingWithInvalidSiblingPid(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 3
	siblingPid := 0
	assert.False(t, tree.IsSibling(pid, siblingPid))
}

func TestIsSiblingWithInvalidPids(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 0
	siblingPid := 0
	assert.False(t, tree.IsSibling(pid, siblingPid))
}

func TestIsSiblingWithSamePids(t *testing.T) {
	// 1 -> 2 -> 3
	//      2 -> 5
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(2, 5)

	pid := 3
	siblingPid := 3
	assert.False(t, tree.IsSibling(pid, siblingPid))
}

func TestIsAncestor(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 9
	ancestorPid := 5
	assert.True(t, tree.IsAncestor(pid, ancestorPid))
}

func TestIsAncestorWithSamePids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 5
	ancestorPid := 5
	assert.False(t, tree.IsAncestor(pid, ancestorPid))
}

func TestIsAncestorWithInvalidDescendantPid(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()

	pid := 0
	ancestorPid := 5
	assert.False(t, tree.IsAncestor(pid, ancestorPid))
}

func TestIsAncestorWithInvalidAncestorPid(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()

	pid := 5
	ancestorPid := 0
	assert.False(t, tree.IsAncestor(pid, ancestorPid))
}

func TestIsAncestorWithInvalidPids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()

	pid := 0
	ancestorPid := 0
	assert.False(t, tree.IsAncestor(pid, ancestorPid))
}

func TestIsDescendant(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 5
	descendantPid := 9
	assert.True(t, tree.IsDescendant(pid, descendantPid))
}

func TestIsDescendantWithSamePids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 5
	descendantPid := 5
	assert.False(t, tree.IsDescendant(pid, descendantPid))
}

func TestIsDescendantWithInvalidPid(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 0
	descendantPid := 9
	assert.False(t, tree.IsDescendant(pid, descendantPid))
}

func TestIsDescendantWithInvalidDescendantPid(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()
	tree.AddProcess(1, 2)
	tree.AddProcess(2, 3)
	tree.AddProcess(3, 4)
	tree.AddProcess(1, 5)
	tree.AddProcess(5, 6)
	tree.AddProcess(6, 7)
	tree.AddProcess(5, 8)
	tree.AddProcess(8, 9)

	pid := 5
	descendantPid := 0
	assert.False(t, tree.IsDescendant(pid, descendantPid))
}

func TestIsDescendantWithInvalidPids(t *testing.T) {
	// 1 -> 2 -> 3 -> 4
	// 1 -> 5 -> 6 -> 7
	//      5 -> 8 -> 9
	tree := NewProcessTree()

	pid := 0
	descendantPid := 0
	assert.False(t, tree.IsDescendant(pid, descendantPid))
}
