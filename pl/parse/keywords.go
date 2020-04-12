package parse

import (
	"shanhu.io/smlvm/lexing"
)

var gKeywords = lexing.KeywordSet(
	"func", "var", "const", "struct", "import", "interface",
	"if", "else", "for", "break", "continue", "return",
	"switch", "case", "default", "fallthrough",
)

var golikeKeywords = lexing.KeywordSet(
	"func", "var", "const", "struct", "import",
	"if", "else", "for",
	"break", "continue", "return",
	"package", "type",
)
