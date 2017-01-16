package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func refAddress(b *builder, opTok *lexing.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	opPos := opTok.Pos

	bref := B.R()
	if types.IsType(bref.T) || !bref.IsSingle() {
		b.CodeErrorf(opPos, "pl.refAdrress.notSingle",
			"%q on %s", op, bref)
		return nil
	} else if !bref.Addressable {
		b.CodeErrorf(opPos, "pl.refAddress.notAddressable",
			"reading address of non-addressable")
		return nil
	}

	r := tast.NewRef(&types.Pointer{T: bref.T})
	return &tast.OpExpr{Op: opTok, B: B, Ref: r}
}

func binaryOpNil(b *builder, opTok *lexing.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		ref := tast.NewRef(types.Bool)
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: ref}
	}

	b.CodeErrorf(opTok.Pos, "pl.invalidExprStmt", "%q on nils", op)
	return nil
}

func binaryOpPtr(b *builder, opTok *lexing.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	atyp := A.R().T
	btyp := B.R().T

	switch op {
	case "==", "!=":
		if types.IsNil(atyp) {
			A = tast.NewCast(A, btyp)
		} else if types.IsNil(btyp) {
			B = tast.NewCast(B, atyp)
		}

		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: tast.NewRef(types.Bool)}
	}

	b.CodeErrorf(opTok.Pos, "pl.invalidExprStmt",
		"%q on pointers", op)
	return nil
}

func binaryOpSlice(b *builder, opTok *lexing.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: tast.NewRef(types.Bool)}
	}
	b.CodeErrorf(opTok.Pos, "pl.invalidExprStmt",
		"%q on slices", op)
	return nil
}
