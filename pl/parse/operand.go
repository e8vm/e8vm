package parse

import (
	"shanhu.io/smlvm/pl/ast"
)

func parseOperand(p *parser) ast.Expr {
	if p.See(Ident) || p.See(Int) || p.See(Float) ||
		p.See(String) || p.See(Char) {
		return ast.NewOperand(p.Shift())
	} else if p.SeeKeyword("this") {
		return ast.NewOperand(p.Shift())
	} else if p.SeeOp("(") {
		lp := p.Shift()
		expr := p.parseExpr()
		rp := p.ExpectOp(")")
		if rp == nil {
			return nil
		}

		return &ast.ParenExpr{
			Lparen: lp,
			Rparen: rp,
			Expr:   expr,
		}
	} else if p.SeeOp("[") {
		t := p.parseType()
		if t == nil {
			return nil
		}
		if !p.SeeOp("{") {
			return t
		}

		ret := new(ast.ArrayLiteral)
		ret.Type = t.(*ast.ArrayTypeExpr)
		ret.Lbrace = p.Shift()
		if !p.SeeOp("}") {
			ret.Exprs = parseExprListClosed(p, "}")
			if p.InError() {
				return nil
			}
		}
		ret.Rbrace = p.ExpectOp("}")
		if p.InError() {
			return nil
		}
		return ret
	} else if p.SeeOp("*") || p.SeeKeyword("func") {
		return p.parseType()
	}

	t := p.Token()
	p.CodeErrorf(t.Pos, "pl.expectOperand",
		"expect an operand, got %s", tokenTypeStr(t))

	return nil
}
