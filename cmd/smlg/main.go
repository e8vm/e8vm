package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl"
)

func handleErrs(errs []*lexing.Error) {
	if errs == nil {
		return
	}
	for _, err := range errs {
		fmt.Println(err)
	}
	os.Exit(-1)
}

func main() {
	pkg := flag.String("pkg", "/...", "package to build")
	homeDir := flag.String("home", ".", "the home directory")
	plan := flag.Bool("plan", false, "plan only")
	golike := flag.Bool("golike", false, "uses go-like syntax")
	runTests := flag.Bool("test", true, "run tests")
	std := flag.String("std", "/std", "stdlib package")
	initPC := flag.Uint("initpc", arch.InitPC, "init PC register value")
	initSP := flag.Uint("initsp", 0, "init SP value, for testing")
	staticOnly := flag.Bool("static", false, "do static analysis only")
	flag.Parse()

	memHome := pl.MakeMemFS()
	dirHome := builds.NewDirFS(*homeDir)
	in := builds.NewOverlay(dirHome, memHome)
	langSet := pl.MakeLangSet(*golike)

	out := builds.NewDirFS(path.Join(*homeDir, "_"))
	b := builds.NewBuilder(in, langSet, *std, out)
	b.Verbose = true
	b.InitPC = uint32(*initPC)
	b.InitSP = uint32(*initSP)
	b.RunTests = *runTests
	b.StaticOnly = *staticOnly

	pkgs, err := builds.SelectPkgs(in, langSet, *pkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !*plan {
		handleErrs(b.BuildPkgs(pkgs))
	} else {
		buildOrder, errs := b.Plan(pkgs)
		handleErrs(errs)
		for _, p := range buildOrder {
			fmt.Println(p)
		}
	}
}
