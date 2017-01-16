package sempass

import (
	"math"
	"strconv"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/parse"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildInt(b *builder, op *lexing.Token) tast.Expr {
	ret, e := strconv.ParseInt(op.Lit, 0, 64)
	if e != nil {
		b.CodeErrorf(op.Pos, "pl.cannotParseConst",
			"invalid integer: %s", e)
		return nil
	}
	// ret can be a negative number? no?
	if ret > math.MaxUint32 {
		b.CodeErrorf(op.Pos, "pl.cannotCast.integerOverFlowed",
			"integer too large to fit in 32-bit")
		return nil
	}

	ref := tast.NewConstRef(types.NewNumber(ret), ret)
	return tast.NewConst(ref)
}

func buildFloat(b *builder, op *lexing.Token) tast.Expr {
	s := op.Lit
	var ret int64
	var index int
	for i, r := range s {
		if r == '.' {
			b.CodeErrorf(
				op.Pos, "pl.notYetSupported",
				"float is not yet supported")
			return nil
		}
		if r == 'e' || r == 'E' {
			index = i
			idxPart := s[0:i]
			idxNum, e := strconv.ParseInt(idxPart, 0, 64)
			if e != nil {
				b.CodeErrorf(op.Pos, "pl.cannotParseConst",
					"invalid integer: %s", e)
				return nil
			}
			ret = idxNum
			break
		}
	}
	expPart := s[index+1:]
	if len(expPart) == 0 {
		b.CodeErrorf(
			op.Pos, "pl.cannotParseConst.wrongFloatFormat",
			"malformed exponent part for the number: %s", op.Lit)
		return nil
	}
	if s[index+1] == 45 { // 45 for '-'
		b.CodeErrorf(
			op.Pos, "pl.notYetSupported",
			"negative exponent index is not yet supported")
		return nil
	}
	expNum, e := strconv.ParseInt(expPart, 0, 64)
	if e != nil {
		b.CodeErrorf(op.Pos, "pl.cannotParseConst",
			"invalid integer: %s", e)
		return nil
	}

	if expNum > 10 {
		b.CodeErrorf(op.Pos, "pl.cannotCast.integerOverFlowed",
			"integer too large to fit in 32-bit")
		return nil
	}
	for ; expNum > 0; expNum-- {
		ret *= 10
	}

	if ret > math.MaxUint32 {
		b.CodeErrorf(op.Pos, "pl.cannotCast.integerOverFlowed",
			"integer too large to fit in 32-bit")
		return nil
	}

	ref := tast.NewConstRef(types.NewNumber(ret), ret)
	return tast.NewConst(ref)
}

func buildChar(b *builder, op *lexing.Token) tast.Expr {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid char: %s", e)
		return nil
	} else if len(v) != 1 {
		b.Errorf(op.Pos, "invalid char in quote: %q", v)
		return nil
	}
	ref := tast.NewConstRef(types.Int8, int64(v[0]))
	return tast.NewConst(ref)
}

func buildString(b *builder, op *lexing.Token) tast.Expr {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid string: %s", e)
		return nil
	}
	ref := tast.NewConstRef(types.String, v)
	return tast.NewConst(ref)
}

func buildIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.CodeErrorf(ident.Pos, "pl.undefinedIdent",
			"undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymVar, tast.SymField:
		ref := tast.NewAddressableRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	case tast.SymConst, tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	case tast.SymFunc:
		if t, ok := t.(*types.Func); ok {
			if t.MethodFunc == nil {
				return &tast.Ident{Token: ident, Ref: tast.NewRef(t), Sym: s}
			}
			if b.this == nil {
				panic("this missing")
			}
			ref := &tast.Ref{T: t.MethodFunc, Recv: b.this}
			return &tast.Ident{Token: ident, Ref: ref, Sym: s}
		}
		return &tast.Ident{Token: ident, Ref: tast.NewRef(t), Sym: s}
	default:
		b.Errorf(ident.Pos, "todo: token type: %s", tast.SymStr(s.Type))
		return nil
	}
}

func buildConstIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.CodeErrorf(ident.Pos, "pl.undefinedIdent",
			"undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymConst:
		ref := tast.NewRef(t)
		return tast.NewConst(ref)
	case tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	}

	b.CodeErrorf(ident.Pos, "pl.expectConst", "%s is a %s; expect a const",
		ident.Lit, tast.SymStr(s.Type),
	)
	return nil
}

func buildConstOperand(b *builder, op *ast.Operand) tast.Expr {
	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Float:
		return buildFloat(b, op.Token)
	case parse.Ident:
		return buildConstIdent(b, op.Token)
	}

	b.Errorf(op.Token.Pos, "expect a constant")
	return nil
}

func buildOperand(b *builder, op *ast.Operand) tast.Expr {
	if op.Token.Type == parse.Keyword && op.Token.Lit == "this" {
		if b.this == nil {
			b.Errorf(op.Token.Pos, "using this out of a method function")
			return nil
		}
		return &tast.This{Ref: b.this}
	}

	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Float:
		return buildFloat(b, op.Token)
	case parse.Char:
		return buildChar(b, op.Token)
	case parse.String:
		return buildString(b, op.Token)
	case parse.Ident:
		return buildIdent(b, op.Token)
	default:
		b.Errorf(op.Token.Pos, "invalid or not implemented: %d",
			op.Token.Type,
		)
		return nil
	}
}
