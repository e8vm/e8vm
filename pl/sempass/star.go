package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildStarExpr(b *builder, expr *ast.StarExpr) tast.Expr {
	hold := b.lhsSwap(false)
	defer b.lhsRestore(hold)

	opPos := expr.Star.Pos

	addr := b.buildExpr(expr.Expr)
	if addr == nil {
		return nil
	}

	addrRef := addr.R()
	if !addrRef.IsSingle() {
		b.CodeErrorf(opPos, "pl.star.onNotSingle", "* on %s", addrRef)
		return nil
	}
	if t, ok := addrRef.T.(*types.Type); ok {
		return tast.NewType(&types.Pointer{T: t.T})
	}

	t, ok := addrRef.T.(*types.Pointer)
	if !ok {
		b.CodeErrorf(opPos, "pl.star.onNonPointer", "* on non-pointer")
		return nil
	}

	r := tast.NewAddressableRef(t.T)
	return &tast.StarExpr{Expr: addr, Ref: r}
}
