package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	logFatal = log.Fatal
	ignore   = []string{".git"}
)

type pair struct {
	hash string
	file string
}
type md5Table map[string][]string

func isIgnored(dir string) bool {
	for _, e := range ignore {
		if e == dir {
			return true
		}
	}
	return false
}

func readTree(directory string) ([]string, error) {
	files := []string{}
	walk := func(path string, fInfo os.DirEntry, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		if fInfo.Type().IsDir() && isIgnored(fInfo.Name()) {
			return filepath.SkipDir
		}

		if fInfo.Type().IsRegular() {
			files = append(files, path)
		}
		return nil
	}

	return files, filepath.WalkDir(directory, walk)
}

func getHash(path string) (pair, error) {
	file, err := os.Open(path)
	if err != nil {
		return pair{}, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return pair{}, err
	}

	ret := pair{
		hash: fmt.Sprintf("%x", hash.Sum(nil)),
		file: path,
	}
	return ret, nil
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
	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	if len(os.Args) < 2 {
		logFatal("Missing required parameter: '<path>'")
	}

	files, err := readTree(os.Args[1])
	if err != nil {
		logFatal(err)
	}

	hashTable := make(md5Table)
	for _, f := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pair, err := getHash(f)
			if err != nil {
				logFatal(err)
			}
			mu.Lock()
			defer mu.Unlock()
			hashTable[pair.hash] = append(hashTable[pair.hash], pair.file)
		}()
	}
	wg.Wait()
	if len(hashTable) == 0 {
		fmt.Println("No duplicate files found.")
	} else {
		showOutput(hashTable)
	}
	os.Exit(0)
}
