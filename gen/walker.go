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
	"strings"

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

		visitor := &errFinder{
			fset: pkg.Fset,
			info: pkg.TypesInfo,
			pkg:  pkg,
		}

		for _, file := range pkg.Syntax {
			ast.Walk(visitor, file)
		}
	}
}

// errFinder is a visitor that looks for grr.Errorf calls and prints information about them.
type errFinder struct {
	fset *token.FileSet
	info *types.Info
	pkg  *packages.Package
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
	nameCounts := map[string]int{}
	imports := utils.NewSetFromSlice(GenDefaultImports())

	for _, arg := range callExpr.Args {
		imp := f.getPackageForExpr(arg)

		if imp != nil {
			imports.Add(imp.Path())
		}

		args = append(args, ExprToArg(f.fset, arg, f.info, nameCounts))
	}

	// pop the first argument, which is the format string
	msg := args[0].Expr
	args = args[1:]

	imports.Add(grrNode.PkgImportPath)

	// generate the error function
	errFunc, err := GenerateErrorFile(
		GenerateFileArgs{
			Args:    args,
			Imports: imports.ToSlice(),
			PkgName: f.pkg.Name,
			ErrMsg:  msg,
		},
	)

	if err != nil {
		fmt.Printf("Error generating error function: %s\n", err)
		grr.Trace(err)
		return nil
	}

	// get an arbitrary Go file from the package
	if len(f.pkg.GoFiles) == 0 {
		fmt.Printf("No Go files found in package: %s\n", f.pkg.PkgPath)
		return nil
	}

	file := f.pkg.GoFiles[0]

	// strip the file name from the path
	pkgPath := filepath.Dir(file)

	writePath := filepath.Join(pkgPath, "grr.gen.go")

	fmt.Printf("Writing to: %s\n", writePath)

	err = os.WriteFile(writePath, []byte(errFunc), 0644)

	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		grr.Trace(err)
		return nil
	}

	return f
}

// getPackageForExpr returns the package for a given identifier, if available
func (f *errFinder) getPackageForExpr(expr ast.Expr) *types.Package {
	switch e := expr.(type) {
	case *ast.Ident:
		// Handle identifiers which may be package names
		if obj, ok := f.info.Uses[e]; ok {
			if pkg, ok := obj.(*types.PkgName); ok {
				return pkg.Imported()
			}
		}
	case *ast.SelectorExpr:
		// Handle qualified identifiers (e.g., pkg.Type)
		if ident, ok := e.X.(*ast.Ident); ok {
			if obj, ok := f.info.Uses[ident]; ok {
				if pkg, ok := obj.(*types.PkgName); ok {
					return pkg.Imported()
				}
			}
		}
	case *ast.CompositeLit:
		// Composite literals like structs or arrays with types defined by SelectorExpr
		if selExpr, ok := e.Type.(*ast.SelectorExpr); ok {
			return f.getPackageForExpr(selExpr)
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
		Pos:           f.fset.Position(callExpr.Pos()),
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

func ExprToArg(fset *token.FileSet, arg ast.Expr, info *types.Info, nameCounts map[string]int) FnArg {
	var buf bytes.Buffer
	var ttype string

	printer.Fprint(&buf, fset, arg)

	if tv, ok := info.Types[arg]; ok && tv.Type != nil {
		ttype = tv.Type.String()

	} else {
		ttype = "any"
	}

	name := generateName(arg, info, nameCounts)

	return FnArg{
		Name: name,
		Expr: buf.String(),
		Type: ttype,
	}
}

func generateName(arg ast.Expr, info *types.Info, nameCounts map[string]int) string {
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

// Simplify type names by extracting the last part and making it camelCase with the package prefix.
func simplifyTypeName(typeName string) string {
	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		// Create camelCase name from package and type
		return strings.ToLower(parts[0]) + strings.Title(parts[1])
	}
	return strings.ToLower(parts[0]) // Just the type name, lowercased
}
