package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/asm"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/dasm"
	"shanhu.io/smlvm/lexing"
)

var (
	doDasm      = flag.Bool("d", false, "do dump")
	ncycle      = flag.Int("n", 100000, "max cycles to execute")
	memSize     = flag.Int("m", 0, "memory size; 0 for full 4GB")
	printStatus = flag.Bool("s", false, "print status after execution")
	randSeed    = flag.Int64("seed", 0, "random seed, 0 for using the time")
)

func run(bs []byte) (int, error) {
	// create a single core machine
	m := arch.NewMachine(&arch.Config{
		MemSize:  uint32(*memSize),
		RandSeed: *randSeed,
	})
	if err := m.LoadImageBytes(bs); err != nil {
		return 0, err
	}

	ret, exp := m.Run(*ncycle)
	if *printStatus {
		m.PrintCoreStatus()
	}

	if exp == nil {
		return ret, nil
	}
	return ret, exp
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("need exactly one input file\n")
	}

	fname := args[0]
	var bs []byte
	f := builds.PathFile(fname)
	var es []*lexing.Error
	if strings.HasSuffix(fname, "_bare.s") {
		bs, es = asm.BuildBareFunc(fname, f)
	} else {
		bs, es = asm.BuildSingleFile(fname, f)
	}

	if len(es) > 0 {
		for _, e := range es {
			fmt.Println(e)
		}
		os.Exit(-1)
		return
	}

	if *doDasm {
		lines := dasm.Dasm(bs, arch.InitPC)
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {
		n, e := run(bs)
		fmt.Printf("(%d cycles)\n", n)
		if e != nil {
			if !arch.IsHalt(e) {
				fmt.Println(e)
			}
		} else {
			fmt.Println("(end of time)")
		}
	}
}
