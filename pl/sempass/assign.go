package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func assign(b *builder, dest, src tast.Expr, op *lexing.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()

	ndest := destRef.Len()
	nsrc := srcRef.Len()
	if ndest != nsrc {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.lengthMismatch",
			"cannot assign(len) %s to %s; length mismatch",
			nsrc, ndest)
		return nil
	}

	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		if !r.Addressable {
			b.CodeErrorf(op.Pos, "pl.cannotAssign.notAddressable",
				"assigning to non-addressable")
			return nil
		}

		destType := r.Type()
		srcType := srcRef.At(i).Type()
		if !types.CanAssign(destType, srcType) {
			b.CodeErrorf(op.Pos, "pl.cannotAssign.typeMismatch",
				"cannot assign %s to %s; type mismatch",
				srcType, destType)
			return nil
		}
	}

	// insert casting if needed
	if srcList, ok := tast.MakeExprList(src); ok {
		newList := tast.NewExprList()
		for i, e := range srcList.Exprs {
			t := e.Type()
			if types.IsNil(t) {
				e = tast.NewCast(e, destRef.At(i).Type())
			} else if v, ok := types.NumConst(t); ok {
				e = constCast(b, nil, v, e, destRef.At(i).Type())
				if e == nil {
					panic("bug")
				}
			}
			newList.Append(e)
		}
		src = newList
	}

	return &tast.AssignStmt{Left: dest, Op: op, Right: src}
}

func parseAssignOp(op string) string {
	opLen := len(op)
	if opLen == 0 {
		panic("invalid assign op")
	}
	return op[:opLen-1]
}

func opAssign(b *builder, dest, src tast.Expr, op *lexing.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()
	if !destRef.IsSingle() || !srcRef.IsSingle() {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.notSingle",
			"cannot assign %s %s %s", destRef, op.Lit, srcRef)
		return nil
	} else if !destRef.Addressable {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.notAddressable",
			"assign to non-addressable")
		return nil
	}

	opLit := parseAssignOp(op.Lit)
	destType := destRef.Type()
	srcType := srcRef.Type()

	if opLit == ">>" || opLit == "<<" {
		if v, ok := types.NumConst(srcType); ok {
			src = constCast(b, op.Pos, v, src, types.Uint)
			if src == nil {
				return nil
			}
			srcRef = src.R()
			srcType = types.Uint
		}

		if !canShift(b, destType, srcType, op.Pos, opLit) {
			return nil
		}
		return &tast.AssignStmt{Left: dest, Op: op, Right: src}
	}

	if v, ok := types.NumConst(srcType); ok {
		src = constCast(b, op.Pos, v, src, destType)
		if src == nil {
			return nil
		}
		srcRef = src.R()
		srcType = destType
	}

	if ok, t := types.SameBasic(destType, srcType); ok {
		switch t {
		case types.Int, types.Int8, types.Uint, types.Uint8:
			return &tast.AssignStmt{Left: dest, Op: op, Right: src}
		}
	}

	b.Errorf(op.Pos, "invalid %s %s %s", destType, opLit, srcType)
	return nil
}

func buildAssignStmt(b *builder, stmt *ast.AssignStmt) tast.Stmt {
	hold := b.lhsSwap(true)
	left := b.buildExpr(stmt.Left)
	b.lhsRestore(hold)
	if left == nil {
		return nil
	}

	right := b.buildExpr(stmt.Right)
	if right == nil {
		return nil
	}

	if stmt.Assign.Lit == "=" {
		return assign(b, left, right, stmt.Assign)
	}

	return opAssign(b, left, right, stmt.Assign)
}
