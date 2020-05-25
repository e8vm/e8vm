package parse

import (
	"shanhu.io/smlvm/lexing"
)

// a pipe that replaces end-lines with semicolons
type semiInserter struct {
	x          lexing.Tokener
	save       *lexing.Token
	insertSemi bool
}

// newSemiInserter creates a new statement lexer that inserts
// semicolons into a token stream.
func newSemiInserter(x lexing.Tokener) *semiInserter {
	return &semiInserter{x: x}
}

func makeSemi(p *lexing.Pos, lit string) *lexing.Token {
	return &lexing.Token{Pos: p, Lit: lit, Type: Semi}
}

// Token returns the next token of lexing
func (sx *semiInserter) Token() *lexing.Token {
	if sx.save != nil {
		t := sx.save
		sx.save = nil
		return t
	}

	for {
		t := sx.x.Token()
		switch t.Type {
		case Semi:
			sx.insertSemi = false
		case Operator:
			switch t.Lit {
			case "}", "]", ")", "++", "--":
				sx.insertSemi = true
			default:
				sx.insertSemi = false
			}
		case lexing.EOF:
			if sx.insertSemi {
				sx.insertSemi = false
				sx.save = t
				return makeSemi(t.Pos, "")
			}
		case Endl:
			if sx.insertSemi {
				sx.insertSemi = false
				return makeSemi(t.Pos, "\n")
			}
			continue // ignore this end line
		case lexing.Comment:
			// do nothing
		default:
			sx.insertSemi = true
		}

		return t
	}
}

// Errs returns the list of lexing errors.
func (sx *semiInserter) Errs() []*lexing.Error {
	return sx.x.Errs()
}
