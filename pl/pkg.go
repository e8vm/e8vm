package pl

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/sempass"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

type pkg struct {
	files map[string]*ast.File

	tops      *syms.Table
	testNames []string
	deps      *dagvis.Graph
}

func newPkg(asts map[string]*ast.File) *pkg {
	ret := new(pkg)
	ret.files = asts

	return ret
}

func (p *pkg) build(b *builder, pinfo *builds.PkgInfo) []*lexing.Error {
	tops, deps, tests, errs := buildPkg(b, p.files, pinfo)
	p.tops = tops
	p.testNames = tests
	p.deps = deps
	return errs
}

func newRand() *rand.Rand {
	var buf [8]byte
	_, err := crand.Read(buf[:])
	var seed int64
	if err == nil {
		seed = int64(binary.LittleEndian.Uint64(buf[:]))
	} else {
		seed = time.Now().UnixNano()
	}
	return rand.New(rand.NewSource(seed))
}

func buildTests(b *builder, tops *syms.Table) (
	testList codegen.Ref, testNames []string,
) {
	tests := listTests(tops)
	n := len(tests)

	if n > 100000 {
		b.CodeErrorf(nil, "pl.tooManyTests", "too many tests in the package")
		return
	}

	if n == 0 {
		return
	}

	rand := newRand()
	perm := rand.Perm(n)

	var irs []*codegen.Func
	var names []string
	for _, index := range perm {
		t := tests[index]
		irs = append(irs, t.ref.IR().(*codegen.Func))
		names = append(names, t.name)
	}
	return b.p.NewTestList(":tests", irs), names
}

func fillConsts(consts []*syms.Symbol) {
	for _, c := range consts {
		name := c.Name()
		t := c.ObjType.(types.T)
		c.Obj = &objConst{name: name, ref: newRef(t, nil)}
	}
}

func fillVars(b *builder, vars []*tast.Define) {
	for _, v := range vars {
		for _, sym := range v.Left {
			t := sym.ObjType.(types.T)
			name := sym.Name()
			ref := newAddressableRef(t, b.newGlobalVar(t, name))
			sym.Obj = &objVar{name: name, ref: ref}
		}
	}
}

func fillFuncAlias(funcs []*tast.FuncAlias) {
	for _, f := range funcs {
		sym := f.Sym
		name := sym.Name()
		t := sym.ObjType.(*types.Func)
		sig := makeFuncSig(t)
		fsym := codegen.NewFuncSym(f.Of.Pkg(), f.Of.Name(), sig)
		f.Sym.Obj = &objFunc{
			name:    name,
			ref:     newRef(t, fsym),
			isAlias: true,
		}
	}
}

func fillFuncs(b *builder, funcs []*tast.Func) {
	for _, f := range funcs {
		name := f.Sym.Name()
		t := f.Sym.ObjType.(*types.Func)
		sig := makeFuncSig(t)
		irFunc := b.p.NewFunc(b.anonyName(name), f.Sym.Pos, sig)
		f.Sym.Obj = &objFunc{name: name, ref: newRef(t, irFunc)}
	}
}

func fillMethods(b *builder, methods []*tast.Func) {
	for _, f := range methods {
		name := f.Sym.Name()
		t := f.Sym.ObjType.(*types.Func)
		s := t.Args[0].T.(*types.Pointer).T.(*types.Struct)

		fullName := fmt.Sprintf("%s:%s", s, name)
		sig := makeFuncSig(t)
		irFunc := b.p.NewFunc(fullName, f.Sym.Pos, sig)
		f.Sym.Obj = &objFunc{
			name:     name,
			ref:      newRef(t, irFunc),
			isMethod: true,
		}
	}
}

func buildFuncs(b *builder, funcs []*tast.Func) {
	for _, f := range funcs {
		obj := f.Sym.Obj.(*objFunc)
		buildFunc(b, f, obj.ref.IR().(*codegen.Func))
	}
}

func buildPkg(
	b *builder, files map[string]*ast.File, pinfo *builds.PkgInfo,
) (
	tops *syms.Table, deps *dagvis.Graph,
	testNames []string, errs []*lexing.Error,
) {
	imports := make(map[string]*builds.Package)
	for as, imp := range pinfo.Import {
		imports[as] = imp.Package
	}

	sp := &sempass.Pkg{
		Path:    b.path,
		Files:   files,
		Imports: imports,
	}

	tops = syms.NewTable()
	b.scope.PushTable(tops)
	defer b.scope.Pop()

	res, depGraph, errs := sp.Build(b.scope)
	if errs != nil {
		return nil, nil, nil, errs
	}

	fillConsts(res.Consts)
	fillVars(b, res.Vars)
	fillFuncAlias(res.FuncAliases)
	fillFuncs(b, res.Funcs)
	fillMethods(b, res.Methods)
	buildFuncs(b, res.Funcs)
	buildFuncs(b, res.Methods)
	addInit(b)
	addStart(b)
	testList, testNames := buildTests(b, tops)
	if testList != nil {
		addTestStart(b, testList, len(testNames))
	}

	return tops, depGraph, testNames, nil
}
