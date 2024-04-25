package gen

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/jackHedaya/grr/grr"
	"golang.org/x/tools/go/packages"
)

func LoadPreviousErrors(pkg string) (map[string]GeneratedError, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}, pkg)

	if err != nil {
		return nil, grr.Errorf("FailedToLoad: failed to load package").AddError(err)
	}

	if len(pkgs) != 1 {
		return nil, grr.Errorf("OneExpected: expected one package, got %d", len(pkgs))
	}

	p := pkgs[0]

	// Create the walker
	walker := &prevWalker{
		fset:       token.NewFileSet(),
		info:       p.TypesInfo,
		pkg:        p,
		prevErrors: map[string]GeneratedError{},
	}

	// Walk the package
	for _, file := range p.Syntax {
		ast.Walk(walker, file)
	}

	return walker.prevErrors, nil
}

type prevWalker struct {
	fset *token.FileSet
	info *types.Info
	pkg  *packages.Package
	// Previous errors found in grr.gen.go files
	prevErrors map[string]GeneratedError
}

// Visit implements the ast.Visitor interface for errFinder.
func (walker *prevWalker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return walker
	}

	// We're looking for structs that are generated by grr
	genDecl, ok := n.(*ast.GenDecl)

	if !ok {
		return walker
	}

	if genDecl.Tok != token.TYPE {
		return walker
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)

		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			continue
		}

		grrErrLookup := walker.pkg.Types.Scope().Lookup("Error")

		if grrErrLookup == nil {
			continue
		}

		grrType := grrErrLookup.Type().Underlying().(*types.Interface)
		comparisonType := walker.info.TypeOf(typeSpec.Name)
		ptr := types.NewPointer(comparisonType)

		if !types.Implements(ptr, grrType) {
			continue
		}

		// Get the name of the struct
		name := typeSpec.Name.Name

		// Get the fields of the struct
		fields := []GrrGenErrorField{}

		for _, field := range structType.Fields.List {
			fieldName := field.Names[0].Name

			fieldType := walker.info.TypeOf(field.Type).String()

			fields = append(fields, GrrGenErrorField{
				Name: fieldName,
				Type: fieldType,
			})
		}

		// Get the generated code

		walker.prevErrors[name] = GeneratedError{
			Name:          name,
			Args:          fields,
			GeneratedCode: "", // doesn't matter because prevErrors are not rewritten, only used to prevent collisions
		}
	}

	return walker
}
