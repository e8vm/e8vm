package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

type parser struct {
	f string
	x lexing.Tokener
	*lexing.Parser

	exprFunc func(p *parser) ast.Expr
	typeFunc func(p *parser) ast.Expr
	stmtFunc func(p *parser) ast.Stmt

	golike bool
}

func makeTokener(f string, r io.Reader, golike bool) lexing.Tokener {
	var x lexing.Tokener = newLexer(f, r)

	x = newSemiInserter(x)
	kw := lexing.NewKeyworder(x)
	kw.Ident = Ident
	kw.Keyword = Keyword
	if !golike {
		kw.Keywords = gKeywords
	} else {
		kw.Keywords = golikeKeywords
	}

	return kw
}

func newParser(f string, r io.Reader, golike bool) (*parser, *lexing.Recorder) {
	ret := new(parser)
	ret.f = f
	ret.golike = golike
	x := makeTokener(f, r, golike)
	rec := lexing.NewRecorder(x)
	ret.x = lexing.NewCommentRemover(rec)
	ret.Parser = lexing.NewParser(ret.x, Types)
	return ret, rec
}

func (p *parser) parseType() ast.Expr {
	if p.typeFunc == nil {
		return nil
	}

	return p.typeFunc(p)
}

func (p *parser) parseExpr() ast.Expr {
	if p.exprFunc == nil {
		return nil
	}

	return p.exprFunc(p)
}

func (p *parser) parseStmt() ast.Expr {
	if p.stmtFunc == nil {
		p.ExpectSemi()
		p.skipErrStmt()
		return nil
	}
	return p.stmtFunc(p)
}

func (p *parser) SeeOp(ops ...string) bool {
	t := p.Token()
	if t.Type != Operator {
		return false
	}
	for _, op := range ops {
		if t.Lit == op {
			return true
		}
	}
	return false
}

func (p *parser) AcceptSemi() *lexing.Token {
	if p.InError() {
		return nil
	}

	t := p.Token()
	if t.Type == Operator && (t.Lit == "}" || t.Lit == ")") {
		return t // fake semicolon by operator
	}

	if t.Type != Semi {
		return nil
	}
	return p.Shift()
}

func (p *parser) SeeSemi() bool {
	t := p.Token()
	if t.Type == Semi {
		return true
	}
	if t.Type == Operator && (t.Lit == "}" || t.Lit == ")") {
		return true
	}
	return false
}

func (p *parser) ExpectSemi() *lexing.Token {
	if p.InError() {
		return nil
	}

	t := p.Token()
	if t.Type == Operator && (t.Lit == "}" || t.Lit == ")") {
		return t // fake semicolon by operator
	}

	if t.Type != Semi {
		p.CodeErrorfHere("pl.missingSemi",
			"expect ';', got %s", tokenTypeStr(t))
		return nil
	}
	return p.Shift()
}

func (p *parser) skipErrStmt() bool {
	if !p.InError() {
		return false
	}

	for {
		t := p.Token()
		if t.Type == Semi || t.Type == lexing.EOF {
			break
		} else if p.SeeOp("}") {
			break
		}
		p.Next()
	}
	if p.See(Semi) {
		p.Next()
	}

	p.BailOut()
	return true
}

func (p *parser) SeeKeyword(kw string) bool {
	return p.SeeLit(Keyword, kw)
}

func (p *parser) ExpectOp(op string) *lexing.Token {
	if p.InError() {
		return nil
	}
	t := p.Token()
	if t.Type != Operator || t.Lit != op {
		p.CodeErrorfHere(
			"pl.expectOp", "expect '%s', got %s", op, tokenTypeStr(t),
		)
		return nil
	}

	return p.Shift()
}

func (p *parser) ExpectKeyword(kw string) *lexing.Token {
	if !p.SeeLit(Keyword, kw) {
		p.CodeErrorfHere(
			"pl.expectKeyword",
			"expect keyword '%s', got %s", kw, tokenTypeStr(p.Token()),
		)
		return nil
	}
	return p.Shift()
}
