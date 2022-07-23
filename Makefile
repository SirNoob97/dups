SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:

clean:
	@rm -rf bin/ testdata/

test: clean
	@./test.sh

.PHONY: build test clean
