package parse

import (
	"shanhu.io/smlvm/lexing"
)

// StmtLexer replaces end-lines with semicolons
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
		case Lbrace, Semi:
			sx.insertSemi = false
		case lexing.EOF:
			if sx.insertSemi {
				sx.insertSemi = false
				sx.save = t
				return &lexing.Token{
					Type: Semi,
					Lit:  t.Lit,
					Pos:  t.Pos,
				}
			}
		case Rbrace:
			if sx.insertSemi {
				sx.save = t
				return &lexing.Token{
					Type: Semi,
					Lit:  ";",
					Pos:  t.Pos,
				}
			}
			sx.insertSemi = true
		case Endl:
			if sx.insertSemi {
				sx.insertSemi = false
				return &lexing.Token{
					Type: Semi,
					Lit:  "\n",
					Pos:  t.Pos,
				}
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
