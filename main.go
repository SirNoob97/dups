package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	logFatal = log.Fatal
	ignore   = []string{".git"}
)

type fileHash struct {
	hash string
	path string
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

func readTree(directory string, paths chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	walk := func(path string, fInfo os.DirEntry, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		if fInfo.Type().IsDir() && isIgnored(fInfo.Name()) {
			return filepath.SkipDir
		}

		if fInfo.Type().IsDir() && directory != path {
			wg.Add(1)
			go readTree(path, paths, wg)
			return filepath.SkipDir
		}

		if fInfo.Type().IsRegular() {
			paths <- path
		}

		return nil
	}

	return filepath.WalkDir(directory, walk)
}

func buildMd5Table(fHash chan fileHash, table chan md5Table) {
	hashTable := make(md5Table)
	for fh := range fHash {
		hashTable[fh.hash] = append(hashTable[fh.hash], fh.path)
	}

	table <- hashTable
}

func getHash(path string) fileHash {
	file, err := os.Open(path)
	if err != nil {
		logFatal(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		logFatal(err)
	}

	ret := fileHash{
		hash: fmt.Sprintf("%x", hash.Sum(nil)),
		path: path,
	}
	return ret
}

func hashFile(paths chan string, fHash chan fileHash, done chan bool) {
	for p := range paths {
		fHash <- getHash(p)

	}
	done <- true
}

func run(directory string) md5Table {
	workers := 2 * runtime.GOMAXPROCS(0)
	paths := make(chan string)
	fHash := make(chan fileHash)
	done := make(chan bool)
	table := make(chan md5Table)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go hashFile(paths, fHash, done)
	}

	go buildMd5Table(fHash, table)

	wg.Add(1)

	err := readTree(directory, paths, wg)
	if err != nil {
		logFatal(err)
	}

	wg.Wait()
	close(paths)

	for i := 0; i < workers; i++ {
		<-done
	}

	close(fHash)

	return <-table
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

	files, err := readTree(os.Args[1])
	if err != nil {
		logFatal(err)
	}

	hashTable := make(md5Table)
	for _, f := range files {
		pair, err := getHash(f)
		if err != nil {
			logFatal(err)
		}

		hashTable[pair.hash] = append(hashTable[pair.hash], pair.file)
	}

	if len(hashTable) == 0 {
		fmt.Println("No duplicate files found.")
	} else {
		showOutput(hashTable)
	}
	os.Exit(0)
}
