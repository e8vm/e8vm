package codegen

import (
	"fmt"

	"shanhu.io/smlvm/link"
)

type strConst struct {
	id   int
	str  string
	pkg  string
	name string
}

func newStrConst(id int, s string) *strConst {
	return &strConst{
		id:  id,
		str: s,
	}
}

func (s *strConst) String() string {
	return fmt.Sprintf("<str %d>", s.id)
}

func (s *strConst) RegSizeAlign() bool { return true }

// A string constant is basically a byte slice.  It contains two register
// size fields: a pointer to the start of the string, and the size of the
// string.
func (s *strConst) Size() int32 {
	return regSize * 2
}

type strPool struct {
	pkg  string
	strs []*strConst
	// strMap is for deduplication
	strMap map[string]*strConst
}

func newStrPool(pkg string) *strPool {
	return &strPool{
		pkg:    pkg,
		strMap: make(map[string]*strConst),
	}
}

func (p *strPool) addString(s string) *strConst {
	exist := p.strMap[s]
	if exist != nil {
		return exist
	}

	n := len(p.strs)
	ret := newStrConst(n, s)
	p.strs = append(p.strs, ret)
	p.strMap[s] = ret
	return ret
}

func countDigit(n int) int {
	ret := 1
	for n > 9 {
		n /= 10
		ret++
	}
	return ret
}

func (p *strPool) declare(lib *link.Pkg) {
	if lib.Path() != p.pkg {
		panic("package name mismatch")
	}

	if len(p.strs) == 0 {
		return
	}

	ndigit := countDigit(len(p.strs))
	nfmt := fmt.Sprintf(":str_%%0%dd", ndigit)

	for i, s := range p.strs {
		s.name = fmt.Sprintf(nfmt, i)
		s.pkg = p.pkg
		v := link.NewVar(0)
		v.Write([]byte(s.str))

		lib.DeclareVar(s.name)
		lib.DefineVar(s.name, v)
	}
}
