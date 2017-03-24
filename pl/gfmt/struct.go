package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printStruct(f *formatter, d *ast.Struct) {
	f.printExprs(d.Kw, " ", d.Name, " ", d.Lbrace)
	f.printEndl()
	f.Tab()
	for i, field := range d.Fields {
		if i != 0 {
			f.printEndPara()
		}
		printIdents(f, field.Idents)
		f.printSpace()
		f.printExprs(field.Type)
	}
	f.printEndl()
	f.ShiftTab()
	f.printToken(d.Rbrace)
}
