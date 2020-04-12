package lexing

// KeywordSet creates a keyword set.
func KeywordSet(words ...string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, k := range words {
		ret[k] = struct{}{}
	}
	return ret
}

// Keyworder contains idents into keywords
type Keyworder struct {
	tokener Tokener

	Keywords map[string]struct{}
	Ident    int
	Keyword  int
}

// NewKeyworder creates a new tokener that changes the type
// of a token into keywords if it is in the keyword map.
func NewKeyworder(tok Tokener) *Keyworder {
	return &Keyworder{tokener: tok}
}

// Token returns the next token, while replacing ident types into
// keyword types if the token is in the keyword set.
func (kw *Keyworder) Token() *Token {
	ret := kw.tokener.Token()
	if kw.Keywords != nil && ret.Type == kw.Ident {
		_, ok := kw.Keywords[ret.Lit]
		if ok {
			ret.Type = kw.Keyword
		}
	}

	return ret
}

// Errs returns the error list on tokening.
func (kw *Keyworder) Errs() []*Error {
	return kw.tokener.Errs()
}
