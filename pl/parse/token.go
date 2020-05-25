package parse

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
)

// The token types used by the parser.
const (
	Keyword = iota
	Ident
	Int
	Float
	Char
	String
	Operator
	Semi
	Endl
)

// Types provides a type name querier
var Types = func() *lexing.Types {
	ret := lexing.NewTypes()
	o := func(t int, name string) {
		ret.Register(t, name)
	}

	o(Keyword, "keyword")
	o(Ident, "identifier")
	o(Int, "integer")
	o(Float, "float")
	o(Char, "char")
	o(String, "string")
	o(Operator, "operator")
	o(Semi, "semicolon")
	o(Endl, "end-line")

	return ret
}()

// TypeStr returns the name of a token type.
func TypeStr(t int) string { return Types.Name(t) }

func tokenTypeStr(t *lexing.Token) string {
	if t.Type == Operator {
		return fmt.Sprintf("'%s'", t.Lit)
	} else if t.Type == Semi {
		return "';'"
	}
	return TypeStr(t.Type)
}
