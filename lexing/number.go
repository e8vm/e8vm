package lexing

// LexNumber lexes a number usign golang's number format.
func LexNumber(x *Lexer, tokInt, tokFloat int) *Token {
	isFloat := false
	start := x.Rune()
	if !IsDigit(start) {
		panic("not starting with a number")
	}

	x.Next()
	r := x.Rune()
	if start == '0' && r == 'x' {
		x.Next()
		for IsHexDigit(x.Rune()) {
			x.Next()
		}
	} else {
		for IsDigit(x.Rune()) {
			x.Next()
		}
		if x.Rune() == '.' {
			isFloat = true
			x.Next()
			for IsDigit(x.Rune()) {
				x.Next()
			}
		}
		if x.Rune() == 'e' || x.Rune() == 'E' {
			isFloat = true
			x.Next()
			if IsDigit(x.Rune()) || x.Rune() == '-' {
				x.Next()
			}
			for IsDigit(x.Rune()) {
				x.Next()
			}
		}
	}
	if isFloat {
		return x.MakeToken(tokFloat)
	}
	return x.MakeToken(tokInt)
}
