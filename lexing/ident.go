package lexing

// IsIdentLetter checks if rune r can start an identifier.
func IsIdentLetter(r rune) bool {
	return r == '_' || IsLetter(r)
}

// LexIdent lexes a typical C/Go langauge identifier.
func LexIdent(x *Lexer, t int) *Token {
	if !IsIdentLetter(x.Rune()) {
		panic("ident must start with letter or _")
	}

	for {
		r, _ := x.Next()
		if !IsIdentLetter(r) && !IsDigit(r) {
			break
		}
	}
	return x.MakeToken(t)
}
