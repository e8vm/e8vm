package pl

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/parse"
)

type importDecl struct {
	as   string
	path string
	pos  *lexing.Pos
}

func listImport(
	f string, o builds.FileOpener, golike bool, lst *builds.ImportList,
) []*lexing.Error {
	rc, err := o.Open()
	if err != nil {
		return lexing.SingleErr(err)
	}
	fast, _, errs := parse.File(f, rc, golike)
	if errs != nil {
		return errs
	}

	if fast.Imports == nil {
		return nil
	}

	m := make(map[string]*importDecl)
	log := lexing.NewErrorList()

	for _, d := range fast.Imports.Decls {
		p, as, e := ast.ImportPathAs(d)
		if e != nil {
			log.CodeErrorf(d.Path.Pos, "pl.invalidImport",
				"invalid path string %s", d.Path.Lit)
			continue
		}

		pos := ast.ImportPos(d)
		if other, found := m[as]; found {
			log.CodeErrorf(pos, "pl.duplImport", "%s already imported", as)
			log.CodeErrorf(other.pos, "pl.duplImport.previousPos",
				"  previously imported here")
			continue
		}

		m[as] = &importDecl{as: as, path: p, pos: pos}
	}

	if errs := log.Errs(); errs != nil {
		return errs
	}

	for as, d := range m {
		lst.Add(as, d.path, d.pos)
	}

	return nil
}
