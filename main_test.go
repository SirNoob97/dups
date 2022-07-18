package main

import (
	"fmt"
	"testing"
	"time"
)

const TEST_DATA = "testdata"

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
		t.Fatal("Expected a non empty Map")
		t.Fatalf("Expected 'nil' got %s", err)
	}
}

// Test_ReadTree_EmptyMapAndNonNil_WhenDirectoryDoesntExists ...
func Test_ReadTree_EmptyMapAndNonNil_WhenDirectoryDoesntExists(t *testing.T) {
	table, err := readTree("")

	if len(table) > 0 && err == nil {
		t.Fatalf("Expected an empty Map got with len of %d", len(table))
		t.Fatal("Expected a non 'nil' error")
	}
}
