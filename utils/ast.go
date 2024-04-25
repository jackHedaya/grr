package utils

import (
	"go/ast"
	"go/types"
)

// getPackageForExpr returns the package for a given identifier, if available
func GetPackageForExpr(info *types.Info, expr ast.Expr) *types.Package {
	switch e := expr.(type) {
	case *ast.Ident:
		// Handle identifiers which may be package names
		if obj, ok := info.Uses[e]; ok {
			if pkg, ok := obj.(*types.PkgName); ok {
				return pkg.Imported()
			}
		}
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			if obj, ok := info.Uses[ident]; ok {
				if pkg, ok := obj.(*types.PkgName); ok {
					return pkg.Imported()
				}
			}
		}
	case *ast.CompositeLit:
		if selExpr, ok := e.Type.(*ast.SelectorExpr); ok {
			return GetPackageForExpr(info, selExpr)
		}
	}

	return nil
}
