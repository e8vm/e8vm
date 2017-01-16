package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseStruct(p *parser) *ast.Struct {
	var ret *ast.Struct
	if !p.golike {
		if p.SeeKeyword("struct") {
			ret = &ast.Struct{
				Kw:     p.ExpectKeyword("struct"),
				Name:   p.Expect(Ident),
				Lbrace: p.ExpectOp("{"),
			}
		} else if p.SeeLit(Ident, "type") {
			ret = &ast.Struct{
				Kw:      p.Expect(Ident),
				Name:    p.Expect(Ident),
				KwAfter: p.ExpectKeyword("struct"),
				Lbrace:  p.ExpectOp("{"),
			}
			p.CodeErrorf(
				ret.Kw.Pos, "pl.invalidStructDecl",
				`G langauge uses "struct %s {}" `+
					`rather than "type %s struct {}"`,
				ret.Name.Lit, ret.Name.Lit,
			)
		} else {
			panic("expect type or struct")
		}
	} else {
		if !p.SeeKeyword("type") {
			panic("expect keyword type")
		}
		ret = &ast.Struct{
			Kw:      p.ExpectKeyword("type"),
			Name:    p.Expect(Ident),
			KwAfter: p.ExpectKeyword("struct"),
			Lbrace:  p.ExpectOp("{"),
		}
	}

	for !p.SeeOp("}") && !p.See(lexing.EOF) {
		if p.SeeKeyword("func") && !p.golike {
			break
		}

		idents := parseIdentList(p)
		if p.skipErrStmt() {
			continue
		}

		field := new(ast.Field)
		field.Idents = idents
		field.Type = p.parseType()

		if p.skipErrStmt() {
			continue
		}

		field.Semi = p.ExpectSemi()
		if p.skipErrStmt() {
			continue
		}

		ret.Fields = append(ret.Fields, field)
	}

	if !p.golike && p.inlineMethod {
		for p.SeeKeyword("func") {
			f := parseFunc(p)
			if f != nil {
				ret.Methods = append(ret.Methods, f)
			}
		}
	}

	ret.Rbrace = p.ExpectOp("}")
	ret.Semi = p.ExpectSemi()

	return ret
}
