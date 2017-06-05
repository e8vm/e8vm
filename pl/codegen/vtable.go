package codegen

import "fmt"

type vtable struct {
	// interface name: pkg.interface
	i string
	// struct name:pkg.structName
	// sName can used for reflection
	s string
	// entries will be sorted by func names
	entries []*FuncSym
}

func newVtable(i, s string) *vtable {
	return &vtable{
		i:       i,
		s:       s,
		entries: make([]*FuncSym, 0),
	}
}

func (t *vtable) String() string {
	return fmt.Sprintf("<vTable %s %s>", t.i, t.s)
}

func (t *vtable) RegSizeAlign() bool { return true }

func (t *vtable) Size() int32 {
	return regSize * 2
}

func (t *vtable) fill(fs []*FuncSym) {
	t.entries = append(t.entries, fs...)
}

type vtablePool struct {
	pkg     string
	vtables []*vtable
}

func newVtablePool(pkg string) *vtablePool {
	return &vtablePool{
		pkg:     pkg,
		vtables: make([]*vtable, 0),
	}
}

func (p *vtablePool) addTable(i, s string) Ref {
	t := newVtable(i, s)
	p.vtables = append(p.vtables, t)
	return t
}
