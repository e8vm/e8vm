package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseElse(p *parser) *ast.ElseStmt {
	if !p.SeeKeyword("else") {
		panic("must start with keyword")
	}

	ret := new(ast.ElseStmt)
	ret.Else = p.Shift()

	if p.SeeKeyword("if") {
		ret.If = p.Shift()
		ret.Expr = parseExpr(p)
	}

	if p.InError() {
		return ret
	}

	if !p.SeeOp("{") {
		p.CodeErrorfHere("pl.parseIf.missingBody",
			"missing if body")
		return ret
	}

	ret.Body = parseBlock(p)
	// might have another else
	if ret.If != nil && p.SeeKeyword("else") {
		ret.Next = parseElse(p)
	}

	return ret
}

func parseIfBody(p *parser) (ret ast.Stmt, isBlock bool) {
	if p.SeeOp("{") {
		return parseBlock(p), true
	}

	if !p.golike {
		if p.SeeKeyword("return") {
			return parseReturnStmt(p, false), false
		} else if p.SeeKeyword("break") {
			return parseBreakStmt(p, false), false
		} else if p.SeeKeyword("continue") {
			return parseContinueStmt(p, false), false
		}
	}

	p.ErrorfHere("expect if body")
	return nil, false
}

// if <cond> { <stmts> }
// if <cond> return <expr>
// if <cond> break
// if <cond> continue
// if <cond> { <stmts> } else { <stmts> }
// if <cond> { <stmts> } else if { <stmts> }
// if <cond> { <stmts> } else if { <stmts> } else { <stmts> }
func parseIfStmt(p *parser) *ast.IfStmt {
	if !p.SeeKeyword("if") {
		panic("must start with keyword")
	}

	ret := new(ast.IfStmt)
	ret.If = p.Shift()
	ret.Expr = parseExpr(p)
	if p.InError() {
		return ret
	}

	var isBlock bool
	ret.Body, isBlock = parseIfBody(p)
	if p.InError() {
		return ret
	}

	if isBlock && p.SeeKeyword("else") {
		// else clause only happens when the body is block
		ret.Else = parseElse(p)
		if p.InError() {
			return ret
		}
	}
	ret.Semi = p.ExpectSemi()
	return ret
}
