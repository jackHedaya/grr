package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"path/filepath"

	"github.com/whiskaway/grr/grr"
	"github.com/whiskaway/grr/utils"
	"golang.org/x/tools/go/packages"
)

var GRR_IMPORT_PATHS = []string{
	"grr/grr",
	"github.com/whiskaway/grr",
	"github.com/whiskaway/grr/grr",
}

var GRR_ERRORF = "Errorf"

// resolveDir ensures the provided path is absolute.
func resolveDir(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving directory: %s\n", err)
		os.Exit(1)
	}
	return absPath
}

// FindAndReplaceErrorsInDir processes all Go files in a directory to find and report grr.Errorf calls.
func FindAndReplaceErrorsInDir(directory string) {
	// Ensure the directory path is absolute
	dir := resolveDir(directory)

	// Set up the configuration to load the packages correctly
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  dir,
	}

	// Load all packages in the directory
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load packages: %s\n", err)
		return
	}

	if len(pkgs) == 0 {
		fmt.Fprintf(os.Stderr, "no packages found in %s\n", dir)
		return
	}

	// Process each package
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			fmt.Fprintf(os.Stderr, "errors in package: %s\n", pkg.Errors[0])
			continue
		}

		visitor := &errFinder{fset: pkg.Fset, info: pkg.TypesInfo}
		for _, file := range pkg.Syntax {
			ast.Walk(visitor, file)
		}
	}
}

// errFinder is a visitor that looks for grr.Errorf calls and prints information about them.
type errFinder struct {
	fset *token.FileSet
	info *types.Info
}

// Visit implements the ast.Visitor interface for errFinder.
func (f *errFinder) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return f
	}

	grrNode, ok := f.getGrrNode(n)

	if !ok {
		return f
	}

	callExpr := grrNode.CallExpr

	pos := f.fset.Position(callExpr.Pos())
	fmt.Printf("Found grr.Errorf call at %s:\n", pos)
	fmt.Println("Arguments and their types:")

	args := []FnArg{}

	for _, arg := range callExpr.Args {
		args = append(args, ExprToArg(f.fset, arg, f.info))
	}

	// pop the first argument, which is the format string
	msg := args[0].Literal
	args = args[1:]

	// generate the error function
	errFunc, err := GenerateErrorFunction(msg, args...)

	if err != nil {
		fmt.Printf("Error generating error function: %s\n", err)
		grr.Trace(err)
		return nil
	}

	os.WriteFile(fmt.Sprintf("grr.gen.go", pos.Filename), []byte(errFunc), 0644)

	return f
}

// getPackageForExpr returns the package for a given identifier, if available
func (f *errFinder) getPackageForExpr(expr *ast.Ident) *types.Package {
	if obj, ok := f.info.Uses[expr]; ok {
		if pkg, ok := obj.(*types.PkgName); ok {
			return pkg.Imported()
		}
	}
	return nil
}

func (f *errFinder) getGrrNode(n ast.Node) (*GrrNode, bool) {
	if n == nil {
		return nil, false
	}

	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	ident, ok := selExpr.X.(*ast.Ident) // Correctly handle type assertion
	if !ok || selExpr.Sel.Name != GRR_ERRORF {
		return nil, false
	}

	pkg := f.getPackageForExpr(ident)
	if pkg == nil || !utils.Contains(GRR_IMPORT_PATHS, pkg.Path()) {
		return nil, false
	}

	return &GrrNode{
		Pos:      f.fset.Position(callExpr.Pos()),
		CallExpr: callExpr,
		SelExpr:  selExpr,
		Ident:    ident,
	}, true
}

type GrrNode struct {
	Pos      token.Position
	CallExpr *ast.CallExpr
	SelExpr  *ast.SelectorExpr
	Ident    *ast.Ident
}

func ExprToArg(fset *token.FileSet, arg ast.Expr, info *types.Info) FnArg {
	var buf bytes.Buffer
	var ttype string

	printer.Fprint(&buf, fset, arg)

	if tv, ok := info.Types[arg]; ok && tv.Type != nil {
		ttype = tv.Type.String()

	} else {
		ttype = "any"
	}

	return FnArg{
		Literal: buf.String(),
		Type:    ttype,
	}
}
