package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	logFatal = log.Fatal
	ignore   = []string{".git"}
)

type md5Table map[string][]string

func isIgnored(dir string) bool {
	for _, e := range ignore {
		if e == dir {
			return true
		}
	}
	return false
}

// TODO: refactor this function, this function should only read a directory tree
// TODO: extract hash operations
func readTree(directory string) (md5Table, error) {
	table := make(md5Table)
	walk := func(path string, fInfo os.DirEntry, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		if fInfo.Type().IsDir() && isIgnored(fInfo.Name()) {
			return filepath.SkipDir
		}

		if fInfo.Type().IsRegular() {
			hash, file := getHash(path)
			table[hash] = append(table[hash], file)
		}
		return nil
	}

	return table, filepath.WalkDir(directory, walk)
}

// TODO: implement a better error handling
func getHash(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		logFatal(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		logFatal(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), path
}

func showOutput(hashTable md5Table) {
	for hash, files := range hashTable {
		if len(files) > 1 {
			fmt.Printf("Files that share the md5 hash: %s\n\n", hash)

			for _, file := range files {
				fmt.Printf("  %s\n", file)
			}
			fmt.Println()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		logFatal("Missing required parameter: '<path>'")
	}

	hashes, err := readTree(os.Args[1])
	if err == nil {
		showOutput(hashes)
	}
}
