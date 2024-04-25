package clean

import (
	"os"
	"path/filepath"

	"github.com/jackHedaya/grr/grr"
	"github.com/jackHedaya/grr/utils"
)

// Deletes all grr.gen.go files in the directory.
func CleanEntry(directory string) error {
	dir, err := utils.ResolveAbsoluteDir(directory)

	if err != nil {
		return grr.Errorf("FailedToResolveDir: Failed to resolve directory").AddError(err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return grr.Errorf("FailedToWalk: Failed to walk %s", path).AddError(err)
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if filepath.Base(path) == "grr.gen.go" {
			err := os.Remove(path)

			if err != nil {
				return grr.Errorf("FailedToDelete: Failed to delete %s", path).AddError(err)
			}
		}

		return nil
	})

	if err != nil {
		return grr.Errorf("FailedToClean: Failed to clean directory").AddError(err)
	}

	return nil
}
