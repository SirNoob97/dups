SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

COMPILER=go
COMPILE_OPTS=build -o
COMPILE=$(COMPILER) $(COMPILE_OPTS)
BUILD_DIR=bin
TEST_DATA_DIR=testdata
test: clean
	@./test.sh

.PHONY: build test clean
