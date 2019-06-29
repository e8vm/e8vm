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
			"cannot assign(len) %d to %d; length mismatch",
			nsrc, ndest)
		return nil
	}

	// check if all addressable
	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		if !r.Addressable {
			b.CodeErrorf(
				op.Pos, "pl.cannotAssign.notAddressable",
				"assigning to non-addressable",
			)
			return nil
		}
	}

	srcTypes := srcRef.TypeList()
	destTypes := destRef.TypeList()
	res := canAssigns(b, op.Pos, destTypes, srcTypes, "assginment")
	if res.err {
		return nil
	}
	if res.needCast {
		src = tast.NewMultiCast(src, destRef, res.castMask)
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
			src = numCast(b, op.Pos, v, src, types.Uint)
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
		src = numCast(b, op.Pos, v, src, destType)
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
