package parse

import (
	"shanhu.io/smlvm/lexing"
)

func lexOperator(x *lexing.Lexer, r rune) *lexing.Token {
	switch r {
	case ';':
		return x.MakeToken(Semi)
	case '{', '}', '(', ')', '[', ']', ',':
		/* do nothing */
	case '/':
		r2 := x.Rune()
		if r2 == '/' || r2 == '*' {
			return lexing.LexComment(x)
		} else if r2 == '=' {
			x.Next()
		}
	case '+', '-', '&', '|':
		r2 := x.Rune()
		if r2 == r || r2 == '=' {
			x.Next()
		}
	case '*', '%', '^', '=', '!', ':':
		r2 := x.Rune()
		if r2 == '=' {
			x.Next()
		}
	case '.':
		r2 := x.Rune()
		if r2 == '.' {
			x.Next()
			r3 := x.Rune()
			if r3 != '.' {
				x.CodeErrorf("pl.invalidDotDot", "expect ..., but see ..")
				return x.MakeToken(Operator)
			}
			x.Next()
		}
	case '>', '<':
		r2 := x.Rune()
		if r2 == r {
			x.Next()
			r3 := x.Rune()
			if r3 == '=' {
				x.Next()
			}
		} else if r2 == '=' {
			x.Next()
		}
	default:
		return nil
	}

	return x.MakeToken(Operator)
}
