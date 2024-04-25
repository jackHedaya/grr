package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
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
		if len(pkg.GoFiles) != len(pkg.Syntax) {
			return grr.Errorf("MismatchedGoFilesAndSyntax: mismatched number of Go files and syntax trees")
		}

		fileToAst := map[string]*ast.File{}

		for idx, astFile := range pkg.Syntax {
			fileToAst[pkg.GoFiles[idx]] = astFile
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

		err = utils.AppendOrCreate(writePath, code)

		if err != nil {
			return grr.Errorf("FailedToWriteFile: failed to write generated file").AddError(err)
		}

		for path, fileNode := range fileToAst {
			err := writeFile(path, pkg.Fset, fileNode)
			if err != nil {
				return grr.Errorf("ASTWriteError: failed to write AST files").AddError(err)
			}
		}
	}

	return nil
}

func writeFile(path string, fset *token.FileSet, fileNode *ast.File) error {
	// backup the original file
	backupPath := path + ".bak"
	err := os.Rename(path, backupPath)

	if err != nil {
		return grr.Errorf("FailedToBackup: failed to backup file").AddError(err)
	}

	restore := func() {
		err := os.Rename(backupPath, path)

		if err != nil {
			fmt.Printf("Failed to restore backup: %v\n", err)
		}
	}

	out, err := os.Create(path)
	if err != nil {
		restore()
		return grr.Errorf("FailedToCreateFile: failed to create file").AddError(err)
	}

	defer out.Close()

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	if err := printer.Fprint(bufWriter, fset, fileNode); err != nil {
		restore()
		return grr.Errorf("FailedToPrint: failed to print AST").AddError(err)
	}

	fmted, err := format.Source(buf.Bytes())

	if err != nil {
		restore()
		return grr.Errorf("FailedToFormat: failed to format source").AddError(err)
	}

	if _, err := out.Write(fmted); err != nil {
		restore()
		return grr.Errorf("FailedToWrite: failed to write to file").AddError(err)
	}

	return nil
}
