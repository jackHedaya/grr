package main

import (
	"fmt"
	"os"

	"github.com/whiskaway/grr/gen"
)

func main() {
	// Read the filename from command line arguments
	if !isProperLength(os.Args) {
		println("Usage: go run main.go <folder.go>")
		os.Exit(1)
	}
	filename := getFilename(os.Args)

	fmt.Printf("Finding and replacing grr.Errorf calls in %s/...\n", filename)

	// Find and replace grr.Errorf calls in the file
	gen.FindAndReplaceErrorsInDir(filename)
}

func isProperLength(args []string) bool {
	// check if either file is passed like ./grr <file> or ./grr -- <file>

	if args[1] == "--" && len(args) == 3 {
		return true
	}

	if len(os.Args) == 2 {
		return true
	}

	return false
}

func getFilename(args []string) string {
	if args[1] == "--" {
		return args[2]
	}
	return args[1]
}
