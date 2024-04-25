package main

import (
	"fmt"
	"os"
	"path"

	"github.com/jackHedaya/grr/clean"
	"github.com/jackHedaya/grr/gen"
	"github.com/jackHedaya/grr/grr"
)

func main() {
	args := os.Args

	if args[1] == "--" {
		// remove the -- from the args
		args = append(args[:1], args[2:]...)
	}

	if len(args) < 3 {
		helpCmd()
		os.Exit(1)
	}

	subCmdArgs := args[2:]

	switch args[1] {
	case "gen":
		genCmd(subCmdArgs)
	case "clean":
		cleanCmd(subCmdArgs)
	case "help":
		helpCmd()
	default:
		helpCmd()
	}
}

// The gen subcommand is used to generate error structs and functions
// subArgs is the arguments passed to the gen subcommand (e.g., ./grr gen <subArgs>)
func genCmd(subArgs []string) {

	if len(subArgs) != 1 {
		fmt.Println("Usage: grr gen <folder>")
		os.Exit(1)
	}

	dirName := subArgs[0]

	if isDir, err := isDir(dirName); !isDir {
		fmt.Println("Usage: grr gen <folder>")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Error checking if path is a directory:", err)
		os.Exit(1)
	}

	fmt.Printf("Finding and replacing grr.Errorf calls in %s/...\n", path.Join(dirName, "..."))
	// Find and replace grr.Errorf calls in the file
	err := gen.GenerateEntry(dirName)

	if err != nil {
		fmt.Printf("Error generating error structs: %s\n", grr.Strace(err))
		os.Exit(1)
	}
}

func cleanCmd(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: grr clean <folder>")
		os.Exit(1)
	}

	dirName := args[0]

	if isDir, err := isDir(dirName); !isDir {
		fmt.Println("Usage: grr gen <folder>")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Error checking if path is a directory:", err)
		os.Exit(1)
	}

	fmt.Printf("Cleaning up grr.Errorf calls in %s\n", path.Join(dirName, "..."))
	err := clean.CleanEntry(dirName)

	if err != nil {
		fmt.Printf("Error cleaning up grr.Errorf calls: %s\n", grr.Strace(err))
		os.Exit(1)
	}
}

func helpCmd() {
	fmt.Println("Usage: grr <command> [<args>]")
	fmt.Println("Commands:")
	fmt.Println("  gen <folder>    Find and replace grr.Errorf calls in the specified folder")
	fmt.Println("  clean <folder>  Clean up grr.Errorf calls in the specified folder")
	fmt.Println("  help            Display this help message")
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, grr.Errorf("PathNotFound: path %s not found", path)
	}

	return info.IsDir(), nil
}
