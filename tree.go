package abnf

import (
	"bytes"
	"strconv"
)

// Tree is tree structure that contains pair of key and value
type Tree struct {
	K       int
	V       []rune
	child   *Tree
	brother *Tree
}

// AllChildren return all children
func (t *Tree) AllChildren() (ret []*Tree) {
	for c := t.child; c != nil; c = c.brother {
		ret = append(ret, c)
	}
	return
}

// Child return first child that Key = k
func (t *Tree) Child(k int) *Tree {
	for c := t.child; c != nil; c = c.brother {
		if c.K == k {
			return c
		}
	}
	return nil
}

// Children return children that Key = k
func (t *Tree) Children(k int) (ret []*Tree) {
	for c := t.child; c != nil; c = c.brother {
		if c.K == k {
			ret = append(ret, c)
		}
	}
	return
}

// HasChild return the tree have child or not
func (t *Tree) HasChild() bool {
	return t.child != nil
}

// add t to p as child
func (t *Tree) add(c *Tree) {
	if ch := t.child; ch == nil {
		t.child = c
	} else {
		for ch.brother != nil {
			ch = ch.brother
		}
		ch.brother = c
	}
	return
}

// GetStack print stack trace
func (t Tree) GetStack() string {
	return t.getStackSub(0)
}

func (t Tree) getStackSub(sp int) string {
	var b bytes.Buffer
	for i := 0; i < sp; i++ {
		b.WriteString("  ")
	}
	b.WriteString("+ ")
	b.WriteString(t.String())
	b.WriteString("\n")
	for c := t.child; c != nil; c = c.brother {
		b.WriteString(c.getStackSub(sp + 1))
	}
	return b.String()
}

func (t Tree) String() string {
	return "[" + strconv.FormatInt(int64(t.K), 10) + "]: " + string(t.V)
}
