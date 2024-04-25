package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
)

type FieldGenerator struct {
	fset       *token.FileSet
	info       *types.Info
	nameCounts map[string]int
}

func NewFieldGenerator(fset *token.FileSet, info *types.Info) *FieldGenerator {
	return &FieldGenerator{
		fset: fset,
		info: info,
		nameCounts: map[string]int{
			// Ensuring that the intrinsic names are unique
			"err":    1,
			"traits": 1,
			"op":     1,
		},
	}
}

func (fg *FieldGenerator) GenerateField(arg ast.Expr) GrrGenErrorField {
	var buf bytes.Buffer
	var ttype string

	printer.Fprint(&buf, fg.fset, arg)

	if tv, ok := fg.info.Types[arg]; ok && tv.Type != nil {
		ttype = tv.Type.String()

	} else {
		ttype = "any"
	}

	name := fg.generateName(arg)

	return GrrGenErrorField{
		Name: name,
		Expr: buf.String(),
		Type: ttype,
	}
}

func (fg *FieldGenerator) generateName(arg ast.Expr) string {
	nameCounts := fg.nameCounts
	info := fg.info

	var name string

	if ident, ok := arg.(*ast.Ident); ok {
		name = ident.Name

	} else if lit, ok := arg.(*ast.BasicLit); ok {
		name = strings.ToLower(lit.Kind.String())

	} else if comp, ok := arg.(*ast.CompositeLit); ok {
		if tv, ok := info.Types[comp]; ok && tv.Type != nil {
			ttype := tv.Type.String()

			switch typ := tv.Type.(type) {
			case *types.Slice:
				elemName := simplifyTypeName(typ.Elem().String())
				name = fmt.Sprintf("%sSlice", elemName)

			case *types.Map:
				keyName := simplifyTypeName(typ.Key().String())
				elemName := simplifyTypeName(typ.Elem().String())
				name = fmt.Sprintf("%s%sMap", keyName, elemName)

			default:
				name = simplifyTypeName(ttype)
			}
		}
	} else {
		name = "arg"
	}

	// Ensuring the name is unique
	if _, ok := nameCounts[name]; ok {
		nameCounts[name]++
		name = fmt.Sprintf("%s%d", name, nameCounts[name])
	} else {
		nameCounts[name] = 1
	}

	return name
}
