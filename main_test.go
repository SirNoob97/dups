package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	TEST_DATA           = "testdata"
	EMPTY_FILE_PATH     = TEST_DATA + "/empty_file"
	EMPTY_FILE_MD5_HASH = "d41d8cd98f00b204e9800998ecf8427e"
)

// Test_IsIgnore_True_WhenDirectoryIsOnTheList
// '.git' director is ignored by default
func Test_IsIgnore_True_WhenDirectoryIsOnTheList(t *testing.T) {
	ignored := isIgnored(".git")

	if !ignored {
		t.Fatal("Expected true got false")
	}
}

// Test_IsIgnore_False_WhenDirectoryIsNotOnTheList
// a random string is used instead a fixed name
func Test_IsIgnore_False_WhenDirectoryIsNotOnTheList(t *testing.T) {
	str := fmt.Sprint(time.Now().UnixMilli())
	ignored := isIgnored(str)

	if ignored {
		t.Fatal("Expected false got true")
	}
}

// Test_IsIgnore_False_WhenParameterIsEmpty ...
func Test_IsIgnore_False_WhenParameterIsEmpty(t *testing.T) {
	ignored := isIgnored("")

	if ignored {
		t.Fatal("Expected false got true")
	}
}

// Test_ReadTree_NilError_WhenTestDataDirectoryHasDuplicatedFiles ...
func Test_ReadTree_NilError_WhenTestDataDirectoryHasDuplicatedFiles(t *testing.T) {
	workers := 2 * runtime.GOMAXPROCS(0)
	paths := make(chan string)
	wg := new(sync.WaitGroup)
	fHash := make(chan fileHash)
	done := make(chan bool)
	table := make(chan md5Table)

	for i := 0; i < workers; i++ {
		go hashFile(paths, fHash, done)
	}

	go buildMd5Table(fHash, table)
	wg.Add(1)

	err := readTree(TEST_DATA, paths, wg)

	wg.Wait()

	close(paths)
	for i := 0; i < workers; i++ {
		<-done
	}
	close(done)

	if err != nil {
		t.Fatalf("Expected 'nil' got %v", err)
	}
}

// Test_ReadTree_NonNilError_WhenDirectoryDoesntExists ...
func Test_ReadTree_NonNilError_WhenDirectoryDoesntExists(t *testing.T) {
	paths := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	err := readTree("", paths, wg)

	wg.Wait()

	if err == nil {
		t.Fatal("Expected a non 'nil' error")
	}
}

// Test_GetHash_FileHashCopy_WhenFileExists ...
func Test_GetHash_FileHashCopy_WhenFileExists(t *testing.T) {
	fHash := getHash(TEST_DATA + "/test_1")

	if len(fHash.hash) == 0 && len(fHash.path) == 0 {
		t.Fatal("Expected a md5 hash and a file path")
	}
}

// Test_GetHash_EmptyFileHashCopy_WhenFileIsAnEmptyFile ...
func Test_GetHash_EmptyFileHashCopy_WhenFileIsAnEmptyFile(t *testing.T) {
	fHash := getHash(EMPTY_FILE_PATH)

	if len(fHash.hash) == 0 && len(fHash.path) == 0 {
		t.Fatal("Expected a md5 hash and a file path")
	}
}

// Test_GetHash_EmptyFileHash_WhenFileIsAnEmptyFile ...
func Test_GetHash_EmptyFileHash_WhenFileIsAnEmptyFile(t *testing.T) {
	fHash := getHash(EMPTY_FILE_PATH)

	if fHash.hash != EMPTY_FILE_MD5_HASH {
		t.Fatalf("Expected %s md5 hash, got %s", EMPTY_FILE_MD5_HASH, fHash)
	}
}

// Test_GetHash_NonNilError_WhenPathIsAnEmptyString ...
func Test_GetHash_NonNilError_WhenPathIsAnEmptyString(t *testing.T) {
	errors := []any{}

	oriLogFata := logFatal
	logFatal = func(v ...any) {
		errors = append(errors, v)
	}

	getHash("")

	logFatal = oriLogFata

	if len(errors) == 0 {
		t.Fatal("Expected a non nil error")
	}
}

// Test_Run_NonEmptyMd5Table_WhenDuplicateFilesAreFound ...
func Test_Run_NonEmptyMd5Table_WhenDuplicateFilesAreFound(t *testing.T) {
	table := run(TEST_DATA)

	if len(table) == 0 {
		t.Fatal("Expected a non empty md5Table")
	}
}

// Test_Run_NonNilError_WhenRootDirectoryIsAnEmptyString ...
func Test_Run_NonNilError_WhenRootDirectoryIsAnEmptyString(t *testing.T)  {
	errors := []any{}

	oriLogFata := logFatal
	logFatal = func(v ...any) {
		errors = append(errors, v)
	}

	run("")

	logFatal = oriLogFata

	if len(errors) == 0 {
		t.Fatal("Expected a non nil error")
	}
}

// Test_ShowOutput_PrintHashTable_WheDuplicateFilesAreFound ...
func Test_ShowOutput_PrintHashTable_WheDuplicateFilesAreFound(t *testing.T) {
	hashTable := make(md5Table)
	hashTable[EMPTY_FILE_MD5_HASH] = []string{EMPTY_FILE_PATH, EMPTY_FILE_PATH}

	tmpFile, err := ioutil.TempFile(TEST_DATA, "temp_file_for_stdout_tests")
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	stdout := os.Stdout
	os.Stdout = tmpFile
	showOutput(hashTable)

	os.Stdout = stdout
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Expected a output message")
	}
}

// Test_ShowOutput_LogError_WhenHashTableIsEmpty ...
func Test_ShowOutput_LogError_WhenHashTableIsEmpty(t *testing.T) {
	hashTable := make(md5Table)
	tmpFile, err := ioutil.TempFile(TEST_DATA, "temp_file_for_stdout_tests")
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	stdout := os.Stdout
	os.Stdout = tmpFile
	showOutput(hashTable)

	os.Stdout = stdout
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	if len(data) != 0 {
		t.Fatal("Expected no output message")
	}
}
