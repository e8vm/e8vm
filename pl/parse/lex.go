package parse

import (
	"io"

	"shanhu.io/smlvm/lexing"
)

func lexG8(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if x.IsWhite(r) {
		panic("incorrect token start")
	}

	switch r {
	case '\n':
		x.Next()
		return x.MakeToken(Endl)
	case '"':
		return lexing.LexString(x, String, '"')
	case '`':
		return lexing.LexRawString(x, String)
	case '\'':
		return lexing.LexString(x, Char, '\'')
	}

	if lexing.IsDigit(r) {
		return lexNumber(x)
	} else if lexing.IsIdentLetter(r) {
		return lexIdent(x)
	}

	// always make progress at this point
	x.Next()
	t := lexOperator(x, r)
	if t != nil {
		return t
	}

	x.CodeErrorf("pl.illegalChar", "illegal char %q", r)
	return x.MakeToken(lexing.Illegal)
}

func newLexer(file string, r io.Reader) *lexing.Lexer {
	return lexing.MakeLexer(file, r, lexG8)
}
