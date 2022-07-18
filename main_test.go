package main

import (
	"fmt"
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

// Test_GetHash_TowNonZeroLengthString_WhenFileExists ...
func Test_GetHash_TowNonZeroLengthString_WhenFileExists(t *testing.T) {
	hash, path := getHash(TEST_DATA + "/test_1")

	if len(hash) == 0 && len(path) == 0 {
		t.Fatal("Expected two string with non zero length")
	}
}

// Test_GetHash_TwoNonZeroLengthString_WhenFileIsAnEmptyFile ...
func Test_GetHash_TwoNonZeroLengthString_WhenFileIsAnEmptyFile(t *testing.T) {
	hash, path := getHash(EMPTY_FILE_PATH)

	if len(hash) == 0 && len(path) == 0 {
		t.Fatal("Expected two string with non zero length")
	}
}

// Test_GetHash_EmptyFileHash_WhenFileIsAnEmptyFile ...
func Test_GetHash_EmptyFileHash_WhenFileIsAnEmptyFile(t *testing.T) {
	hash, _ := getHash(EMPTY_FILE_PATH)

	if hash != EMPTY_FILE_MD5_HASH {
		t.Fatalf("Expected %s md5 hash, got %s", EMPTY_FILE_MD5_HASH, hash)
	}
}

// Test_GetHash_LogErrors_WhenPathIsAnEmptyString ...
func Test_GetHash_LogErrors_WhenPathIsAnEmptyString(t *testing.T) {
	origLogFatal := logFatal
	defer func() {
		logFatal = origLogFatal
	}()

	errors := []any{}
	logFatal = func(a ...any) {
		errors = append(errors, a)
	}

	getHash("")

	if len(errors) == 0 {
		t.Fatal("Expected errors to be logged")
	}
}
