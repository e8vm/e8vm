package lexing

import (
	"io"
)

// LexFunc is a function type that takes a lexer and returns the next token.
type LexFunc func(x *Lexer) *Token

// Lexer parses a file input stream into tokens.
type Lexer struct {
	s *lexScanner

	e    error
	errs *ErrorList

	r rune

	IsWhite WhiteFunc
	LexFunc LexFunc
}

// MakeLexer creates a lexer with the particular lexer func.
func MakeLexer(file string, r io.Reader, f LexFunc) *Lexer {
	ret := NewLexer(file, r)
	ret.LexFunc = f
	return ret
}

// NewLexer creates a new lexer.
func NewLexer(file string, r io.Reader) *Lexer {
	ret := new(Lexer)
	ret.s = newLexScanner(file, r)
	ret.errs = NewErrorList()

	ret.IsWhite = IsWhite
	ret.Next()

	return ret
}

// Next pushes the current rune into the scanning buffer,
// and returns the next rune.
func (x *Lexer) Next() (rune, error) {
	x.r, x.e = x.s.next()
	return x.r, x.e
}

// Rune returns the current rune.
func (x *Lexer) Rune() rune { return x.r }

// See returns true when the current rune is r.
func (x *Lexer) See(r rune) bool { return x.r == r }

// Discard clears the scanning buffer
func (x *Lexer) Discard() { x.s.accept() }

// Ended returns true when the scanning stops.
func (x *Lexer) Ended() bool { return x.e != nil }

// SkipWhite is a helper function that skips
// any rune that returns true by IsWhite function.
// The buffer is discarded after the skipping.
func (x *Lexer) SkipWhite() {
	for {
		if x.Ended() || !x.IsWhite(x.r) {
			break
		}
		x.Next()
	}
	x.Discard()
}

// Buffered returns the current buffered string
// in the scanner
func (x *Lexer) Buffered() string { return x.s.buffered() }

// MakeToken accepts the runes in the scanning buffer
// and returns it as a token of type t.
func (x *Lexer) MakeToken(t int) *Token {
	s, p := x.s.accept()
	return &Token{t, s, p}
}

// Token returns the next parsed token.
// It ends with a token with type EOF.
func (x *Lexer) Token() *Token {
	x.SkipWhite()

	if x.Ended() {
		return x.MakeToken(EOF)
	}

	if x.LexFunc == nil {
		x.Next()
		return x.MakeToken(Illegal)
	}

	return x.LexFunc(x)
}

// Errorf adds an error into the error list with current postion.
func (x *Lexer) Errorf(f string, args ...interface{}) {
	x.errs.CodeErrorf(x.s.startPos(), "", f, args...)
}

// CodeErrorf adds an error into the error list with error code.
func (x *Lexer) CodeErrorf(c, f string, args ...interface{}) {
	x.errs.CodeErrorf(x.s.startPos(), c, f, args...)
}

// Errs returns the lexing errors.
func (x *Lexer) Errs() []*Error {
	if x.e != nil && x.e != io.EOF {
		return []*Error{{Err: x.e}}
	}

	return x.errs.Errs()
}
