package gen

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackHedaya/grr/grr"
	"github.com/jackHedaya/grr/utils"
	"golang.org/x/tools/go/packages"
)

// GenerateEntry processes all Go files in a directory to find and report grr.Errorf calls.
func GenerateEntry(directory string) error {
	// Ensure the directory path is absolute
	dir, err := utils.ResolveAbsoluteDir(directory)

	if err != nil {
		return grr.Errorf("unable to determine directory").AddError(err)
	}

	// Set up the configuration to load the packages correctly
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  dir,
	}

	// Load all packages in the directory
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return grr.Errorf("FailedToLoadPackages: failed to load packages").AddError(err)
	}

	if len(pkgs) == 0 {
		return grr.Errorf("NoPackagesFound: no packages found in directory. string builder for testing: %v", strings.Builder{})
	}

	// Process each package
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			fmt.Printf("Ignoring package %s due to errors: %v\n", pkg.PkgPath, pkg.Errors)
			continue
		}

		// get previous errors
		prevErrors, err := LoadPreviousErrors(pkg.PkgPath)

		if err != nil {
			return grr.Errorf("FailedToLoadPreviousErrors: failed to load previous errors").AddError(err)
		}

		pkgWalker := &grrWalker{
			fset:            pkg.Fset,
			info:            pkg.TypesInfo,
			pkg:             pkg,
			generatedErrors: map[string]GeneratedError{},
			prevErrors:      prevErrors,
			imports:         utils.NewSetFromSlice(GenDefaultImports()),
		}

		// Load previous errors from grr.gen.go files
		// err := visitor.LoadPreviousErrors()

		for _, astFile := range pkg.Syntax {
			ast.Walk(pkgWalker, astFile)
		}

		pkgPath, err := utils.GetPackagePath(pkg)

		if err != nil {
			return grr.Errorf("FailedToGetPackagePath: failed to get package path").AddError(err)
		}

		code, err := GenerateErrorFile(pkg.Name, pkgWalker.imports.ToSlice(), pkgWalker.generatedErrors)

		if err != nil {
			return grr.Errorf("GenerateErrorFile: failed to generate error file").AddError(err)
		}

		if len(pkg.GoFiles) == 0 {
			fmt.Printf("No Go files found in package: %s\n", pkg.PkgPath)
			continue
		}

		if len(pkgWalker.generatedErrors) == 0 {
			fmt.Printf("No grr.Errorf calls found in package: %s\n", pkg.PkgPath)
			continue
		}

		writePath := filepath.Join(pkgPath, "grr.gen.go")

		fmt.Printf("Writing to: %s\n", writePath)

		err = os.WriteFile(writePath, code, fs.FileMode(os.O_APPEND)|fs.FileMode(os.O_CREATE))

		if err != nil {
			return grr.Errorf("FailedToWriteFile: failed to write generated file").AddError(err)
		}

		err = writeFiles(pkg.Fset, pkg, pkgPath)

		if err != nil {
			return grr.Errorf("ASTWriteError: failed to write AST files").AddError(err)
		}
	}

	return nil
}

func writeFiles(fset *token.FileSet, pkg *packages.Package, outputDir string) error {

	for _, file := range pkg.GoFiles {
		outputFile := outputDir + "/" + file

		// Create file
		out, err := os.Create(outputFile)
		if err != nil {
			return err
		}

		// Print the AST back to source code
		if err := printer.Fprint(out, fset, file); err != nil {
			out.Close() // Close file on error
			return err
		}

		// Close the file
		if err := out.Close(); err != nil {
			return err
		}
	}

	return nil
}
