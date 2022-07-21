package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

// Test_ReadTree_NonEmptyMapAndNil_WhenTestDataDirectoryHasDuplicatedFiles ...
func Test_ReadTree_NonEmptyMapAndNil_WhenTestDataDirectoryHasDuplicatedFiles(t *testing.T) {
	table, err := readTree(TEST_DATA)

	if len(table) == 0 && err != nil {
		t.Error("Expected a non empty Map")
		t.Errorf("Expected 'nil' got %v", err)
	}
}

// Test_ReadTree_EmptyMapAndNonNil_WhenDirectoryDoesntExists ...
func Test_ReadTree_EmptyMapAndNonNil_WhenDirectoryDoesntExists(t *testing.T) {
	table, err := readTree("")

	if len(table) > 0 && err == nil {
		t.Errorf("Expected an empty Map got with len of %d", len(table))
		t.Error("Expected a non 'nil' error")
	}
}

// Test_GetHash_PairOfMD5HashWithTheFilePathAndNilError_WhenFileExists ...
func Test_GetHash_PairOfMD5HashWithTheFilePathAndNilError_WhenFileExists(t *testing.T) {
	fHash := getHash(TEST_DATA + "/test_1")

	if len(fHash.hash) == 0 && len(fHash.path) == 0 {
		t.Fatal("Expected a md5 hash and a file path")
	}
}

// Test_GetHash_PairOfMD5HashWithTheFileAndNilError_WhenFileIsAnEmptyFile ...
func Test_GetHash_PairOfMD5HashWithTheFileAndNilError_WhenFileIsAnEmptyFile(t *testing.T) {
	fHash := getHash(EMPTY_FILE_PATH)

	if len(fHash.hash) == 0 && len(fHash.path) == 0 {
		t.Fatal("Expected a md5 hash and a file path")
	}
}

// Test_GetHash_EmptyFileHashAndNilError_WhenFileIsAnEmptyFile ...
func Test_GetHash_EmptyFileHashAndNilError_WhenFileIsAnEmptyFile(t *testing.T) {
	fHash := getHash(EMPTY_FILE_PATH)

	if fHash.hash != EMPTY_FILE_MD5_HASH {
		t.Fatalf("Expected %s md5 hash, got %s", EMPTY_FILE_MD5_HASH, fHash)
	}
}

// Test_GetHash_NonNilError_WhenPathIsAnEmptyString ...
func Test_GetHash_NonNilError_WhenPathIsAnEmptyString(t *testing.T) {
	_, err := getHash("")

	if err == nil {
		t.Fatal("Expected a non nil error")
	}
}

// Test_ShowOutput_PrintHashTable_WheDuplicateFilesAreFound ...
func Test_ShowOutput_PrintHashTable_WheDuplicateFilesAreFound(t *testing.T) {
	files, err := readTree(TEST_DATA)
	if err != nil {
		t.Fatalf("Expected a nil error, got %v", err)
	}

	hashTable := make(md5Table)
	for _, f := range files {
		fHash := getHash(f)
		hashTable[fHash.hash] = append(hashTable[fHash.hash], fHash.path)
	}

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
