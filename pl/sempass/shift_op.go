package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/types"
)

func canShift(b *builder, atyp, btyp types.T, pos *lexing.Pos, op string) bool {
	if !types.IsInteger(atyp) {
		b.CodeErrorf(pos, "pl.cannotShift",
			"cannot %q on %s, require uint", op, atyp)
		return false
	} else if !types.IsInteger(btyp) {
		b.CodeErrorf(pos, "pl.cannotShift",
			"%q with %s", op, btyp)
		return false
	} else if !types.IsUnsigned(btyp) {
		b.CodeErrorf(pos, "pl.cannotShift",
			"%q with %s; must be unsigned", op, btyp)
		return false
	}
	return true
}
