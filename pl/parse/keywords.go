package parse

func keywordSet(s ...string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, k := range s {
		ret[k] = struct{}{}
	}
	return ret
}

var gKeywords = keywordSet(
	"func", "var", "const", "struct", "import",
	"if", "else", "for", "break", "continue", "return",
	"switch", "case", "default", "fallthrough",
)

var golikeKeywords = keywordSet(
	"func", "var", "const", "struct", "import",
	"if", "else", "for",
	"break", "continue", "return",
	"package", "type",
)
