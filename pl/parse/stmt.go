package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseStmt(p *parser) ast.Stmt {
	first := p.Token()
	if first.Type == Keyword {
		switch first.Lit {
		case "const":
			return parseConstDecls(p)
		case "var":
			return parseVarDecls(p)
		case "if":
			return parseIfStmt(p)
		case "for":
			return parseForStmt(p)
		case "switch":
			return parseSwitchStmt(p)
		case "return":
			return parseReturnStmt(p, true)
		case "break":
			return parseBreakStmt(p, true)
		case "continue":
			return parseContinueStmt(p, true)
		case "fallthrough":
			return parseFallthrough(p)
		case "else":
			// a common error case where else leads a statement.
			p.CodeErrorfHere(
				"pl.elseStart",
				`else must be after the if statement,
					and on the same line as the last '}'`,
			)
			p.skipErrStmt()
			return nil
		case "case":
			return swtichErr(p, first.Lit)
		case "default":
			return swtichErr(p, first.Lit)
		}
	}

	if p.SeeOp("{") {
		ret := new(ast.BlockStmt)
		ret.Block = parseBlock(p)
		ret.Semi = p.ExpectSemi()
		return ret
	}

	return parseSimpleStmt(p)
}

func swtichErr(p *parser, s string) ast.Stmt {
	p.CodeErrorfHere(
		"pl.missingSwitch",
		"%s must be within a 'switch' block", s)
	p.skipErrStmt()
	return nil
}

func makeParser(f string, r io.Reader, golike bool) (
	*parser, *lexing.Recorder,
) {
	p, rec := newParser(f, r, golike)
	p.exprFunc = parseExpr
	p.stmtFunc = parseStmt
	p.typeFunc = parseType
	return p, rec
}

// Stmts parses a file input stream as a list of statements,
// like a bare function body.
func Stmts(f string, r io.Reader) ([]ast.Stmt, []*lexing.Error) {
	p, _ := makeParser(f, r, false)

	var ret []ast.Stmt
	for !p.See(lexing.EOF) {
		if stmt := p.parseStmt(); stmt != nil {
			ret = append(ret, stmt)
		}
		p.skipErrStmt()
	}

	if es := p.Errs(); es != nil {
		return nil, es
	}

	return ret, nil
}
