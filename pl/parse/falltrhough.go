package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseFallthrough(p *parser) *ast.FallthroughStmt {
	ret := new(ast.FallthroughStmt)
	ret.Kw = p.ExpectKeyword("fallthrough")
	return ret
}
