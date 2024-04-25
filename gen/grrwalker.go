package gen

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/jackHedaya/grr/grr"
	"github.com/jackHedaya/grr/utils"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

var GRR_IMPORT_PATHS = []string{
	"grr/grr",
	"github.com/jackHedaya/grr",
	"github.com/jackHedaya/grr/grr",
}

var GRR_ERRORF = "Errorf"

// grrWalker is a visitor that looks for grr.Errorf calls and prints information about them.
type grrWalker struct {
	fset *token.FileSet
	info *types.Info
	pkg  *packages.Package
	// A map of error names to their corresponding generated error structs
	generatedErrors map[string]GeneratedError
	// Previous errors found in grr.gen.go files
	prevErrors map[string]GeneratedError
	imports    *utils.Set[string]
}

type GeneratedError struct {
	Name          string
	Args          []GrrGenErrorField
	Msg           string
	GeneratedCode string
}

// Visit implements the ast.Visitor interface for errFinder.
func (walker *grrWalker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return walker
	}

	grrNode, ok := getGrrNode(walker.fset, walker.info, n)

	if !ok {
		return walker
	}

	callExpr := grrNode.CallExpr

	pos := walker.fset.Position(callExpr.Pos())
	fmt.Printf("Found grr.Errorf call at %s:\n", pos)

	args := []GrrGenErrorField{}

	fieldGen := NewFieldGenerator(walker.fset, walker.info)

	for _, arg := range callExpr.Args {
		imp := utils.GetPackageForExpr(walker.info, arg)

		if imp != nil {
			walker.imports.Add(imp.Path())
		}

		args = append(args, fieldGen.GenerateField(arg))
	}

	// pop the first argument, which is the format string
	msg := args[0].Expr
	args = args[1:]

	walker.imports.Add(grrNode.PkgImportPath)

	// generate the error function
	genErr, err := walker.GenerateErrorStruct(
		GenerateFileArgs{
			Args:   args,
			ErrMsg: msg,
		},
	)

	// if _, ok := err.(*ErrNoErrorName); ok {
	// 	fmt.Printf("Error found with no error name in error message: %s. Skipping...\n", msg)
	// 	return walker
	// }

	if err != nil {
		grr.Errorf("FailedToGenerateStruct: failed to generate error struct").AddError(err).Trace()
		return nil
	}

	newFuncIdent := ast.NewIdent(genErr.Name)

	newCallExpr := &ast.CallExpr{
		Fun:  newFuncIdent,
		Args: callExpr.Args,
	}

	// replace the grr.Errorf call with the generated error struct
	astutil.Apply(n, func(cursor *astutil.Cursor) bool {
		if cursor.Node() == callExpr {
			cursor.Replace(newCallExpr)
		}

		return true
	}, nil)

	walker.generatedErrors[genErr.Name] = *genErr

	return walker
}

func getGrrNode(fset *token.FileSet, typesInfo *types.Info, node ast.Node) (*GrrNode, bool) {
	if node == nil {
		return nil, false
	}

	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok || selExpr.Sel.Name != GRR_ERRORF {
		return nil, false
	}

	pkg := utils.GetPackageForExpr(typesInfo, ident)
	if pkg == nil || !utils.Contains(GRR_IMPORT_PATHS, pkg.Path()) {
		return nil, false
	}

	return &GrrNode{
		Pos:           fset.Position(callExpr.Pos()),
		CallExpr:      callExpr,
		SelExpr:       selExpr,
		Ident:         ident,
		PkgImportPath: pkg.Path(),
	}, true
}

type GrrNode struct {
	Pos           token.Position
	CallExpr      *ast.CallExpr
	SelExpr       *ast.SelectorExpr
	Ident         *ast.Ident
	PkgImportPath string
}

// Simplify type names by extracting the last part and making it camelCase with the package prefix.
func simplifyTypeName(typeName string) string {
	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[0]) + strings.Title(parts[1])
	}
	return strings.ToLower(parts[0])
}
