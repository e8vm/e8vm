package link8

type tracer struct {
	lnk  *linker
	hits [][]bool
}

func newTracer(lnk *linker) *tracer {
	npkg := lnk.npkg()
	ret := new(tracer)
	ret.lnk = lnk
	ret.hits = make([][]bool, npkg)
	for i := 0; i < npkg; i++ {
		p := lnk.pkg(i)
		ret.hits[i] = make([]bool, len(p.symbols))
	}

	return ret
}

func (t *tracer) hit(pkg *Pkg, sym uint32) bool {
	i := t.lnk.pkgIndex(pkg.path)
	pt := &t.hits[i][sym]
	ret := *pt
	*pt = true
	return ret
}

// traceUsed traces symbols/objects that are used.
// only these objects need to be linked into the final result.
func traceUsed(lnk *linker, p *Pkg, roots []uint32) []pkgSym {
	t := newTracer(lnk)

	var cur []pkgSym
	for _, index := range roots {
		cur = append(cur, pkgSym{p, index})
	}

	var next []pkgSym
	var ret []pkgSym

	addLink := func(ps pkgSym, lnk *link) {
		pkg := ps.Import(lnk.pkg)
		if t.hit(pkg, lnk.sym) {
			return
		}

		item := pkgSym{pkg, lnk.sym}
		next = append(next, item)
	}

	// BFS traverse all the symbols used by the symbol
	for len(cur) > 0 {
		for _, ps := range cur {
			ret = append(ret, ps)

			typ := ps.Type()
			switch typ {
			case SymFunc:
				f := ps.Func()
				for _, link := range f.links {
					addLink(ps, link)
				}
			case SymVar:
				v := ps.Var()
				for _, link := range v.links {
					addLink(ps, link)
				}
			}
		}

		cur = next
		next = nil
	}

	return ret
}
