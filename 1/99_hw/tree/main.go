package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	const (
		emptyFilePostfix    = "(empty)"
		vertical            = "│\t"
		subcatalogPrefix    = "├───"
		endSubcatalogPrefix = "└───"
	)

	var walk func(relPath string, prefix string)

	walk = func(relPath string, prefix string) {
		entries, err := os.ReadDir(relPath)
		if err != nil {
			fmt.Fprintln(out, err)
			return
		}

		var cleared []fs.DirEntry
		if printFiles {
			cleared = entries
		} else {
			for _, entry := range entries {
				if entry.IsDir() {
					cleared = append(cleared, entry)
				}
			}
		}

		for i, entry := range cleared {
			var (
				currentPrefix string
				nextPrefix    string
			)

			if i == len(cleared)-1 {
				nextPrefix = prefix + "\t"
				currentPrefix = prefix + endSubcatalogPrefix
			} else {
				nextPrefix = prefix + vertical
				currentPrefix = prefix + subcatalogPrefix
			}

			if entry.IsDir() {
				fmt.Fprintf(out, "%s%s\n", currentPrefix, entry.Name())
				walk(filepath.Join(relPath, entry.Name()), nextPrefix)
			} else {
				fileInfo, err := entry.Info()
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				size := emptyFilePostfix
				if fileInfo.Size() > 0 {
					size = fmt.Sprintf("(%db)", fileInfo.Size())
				}
				fmt.Fprintf(out, "%s%s %s\n", currentPrefix, entry.Name(), size)
			}
		}
	}

	walk(path, "")
	return nil
}
