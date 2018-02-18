package builds

import (
	"fmt"
	"io"
	"path"

	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/lexing"
)

func deps(node *dagvis.MapNode) []string {
	depNodes := dagvis.AllInsSorted(node)
	ret := make([]string, 0, len(depNodes))
	for _, dep := range depNodes {
		ret = append(ret, dep.Name)
	}
	return ret
}

func fillImports(c *context, p *pkg) {
	for _, imp := range p.imports {
		imp.Package = c.pkgs[imp.Path].pkg
		if imp.Package == nil {
			panic("bug")
		}
	}
}

func buildMain(c *context, p *pkg) []*lexing.Error {
	lib := p.pkg.Lib
	main := p.pkg.Main

	if main == "" || !lib.HasFunc(main) {
		return nil
	}

	log := lexing.NewErrorList()

	fout, err := c.res.bin(p.path)
	if err != nil {
		return lexing.SingleErr(err)
	}
	lexing.LogError(log, linkPkg(c, fout, p, main))
	lexing.LogError(log, fout.Close())

	return log.Errs()
}

func parseOutput(c *context, p string) func(f string, toks []*lexing.Token) {
	if c.SaveFileTokens == nil {
		return nil
	}
	return func(file string, tokens []*lexing.Token) {
		c.SaveFileTokens(path.Join(p, file), tokens)
	}
}

func makePkgInfo(c *context, p *pkg) *PkgInfo {
	return &PkgInfo{
		Path:   p.path,
		Files:  p.fileSet(),
		Import: p.imports,

		Flags: &Flags{StaticOnly: c.StaticOnly},
		Output: func(name string) (io.WriteCloser, error) {
			return c.res.output(p.path, name)
		},
		ParseOutput: parseOutput(c, p.path),
		AddFuncDebug: func(name string, pos *lexing.Pos, frameSize uint32) {
			c.debugFuncs.Add(p.path, name, pos, frameSize)
		},
	}
}

func buildPkg(c *context, pkg *pkg) []*lexing.Error {
	fillImports(c, pkg)

	compiled, es := pkg.lang.Compile(makePkgInfo(c, pkg))
	if es != nil {
		return es
	}
	pkg.pkg = compiled
	c.linkPkgs[pkg.path] = pkg.pkg.Lib // add for linking

	if c.StaticOnly { // static analysis stops here
		return nil
	}

	if es := buildMain(c, pkg); es != nil {
		return es
	}
	if !pkg.runTests { // skip running tests
		return nil
	}

	return runPkgTests(c, pkg)
}

func build(c *context, pkgs []string) []*lexing.Error {
	for _, p := range pkgs {
		p = c.importPath(p)
		if pkg, es := prepare(c, p); es != nil {
			return es
		} else if pkg.err != nil {
			return lexing.SingleErr(pkg.err)
		}
	}

	if c.RunTests {
		for _, p := range pkgs {
			c.pkgs[p].runTests = true
		}
	}

	g := dagvis.NewGraph(c.deps).Reverse()
	if c.SaveDeps != nil {
		c.SaveDeps(g)
	}

	m, err := dagvis.NewMap(g)
	if err != nil {
		return lexing.SingleErr(err)
	}
	nodes := m.SortedNodes()
	for _, node := range nodes {
		name := node.Name
		if c.Verbose { // report progress
			logln(c, name)
		}

		pkg := c.pkgs[name]
		if pkg == nil {
			panic(fmt.Sprintf("package not found: %q", name))
		}
		pkg.deps = deps(node)
		if es := buildPkg(c, pkg); es != nil {
			return es
		}
	}

	return nil
}
