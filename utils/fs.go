package utils

import (
	"os"
	"path/filepath"

	"github.com/jackHedaya/grr/grr"
	"golang.org/x/tools/go/packages"
)

// resolveDir ensures the provided path is absolute.
func ResolveAbsoluteDir(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", grr.Errorf("FailedToResolveDir: failed to resolve absolute path").AddError(err)
	}
	return absPath, nil
}

func GetPackagePath(pkg *packages.Package) (string, error) {
	file := pkg.GoFiles[0]

	checkPath := filepath.Dir(file)

	if info, err := os.Stat(checkPath); os.IsNotExist(err) || !info.IsDir() {
		return "", grr.Errorf("FailedToGetPackagePath: failed to get package path").AddError(err)
	}

	return checkPath, nil
}
