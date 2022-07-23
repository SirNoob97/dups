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

# This will stay here for quick access in case of modification
TEST_DATA_DIR=testdata
TEST_TREE=$(TEST_DATA_DIR)/{1..5}/{a..e}
TEST_FILES=$(TEST_DATA_DIR)/{test_{1..5},empty_file,{1..5}/{test_{1..5},{a..e}/test_{a..e}}}
FIRST_TEST_MSG=echo "this is a test file" | tee $(TEST_DATA_DIR)/{test_{1,2},{1..5}/{test_{1,2},{a..e}/test_{a,b}}} 1>/dev/null
SECOND_TEST_MSG=echo "this is a another test file" | tee $(TEST_DATA_DIR)/{test_{4,5},{1..5}/{test_{4,5},{a..e}/test_{d,e}}} 1>/dev/null

all: build

$(BUILD_DIR)/dups: clean-$(BUILD_DIR)
> @mkdir -p $(@D)
> @$(COMPILE) $@

build: $(BUILD_DIR)/dups

clean-$(BUILD_DIR):
> @rm -rf $(BUILD_DIR)

clean-$(TEST_DATA_DIR):
> @rm -rf $(TEST_DATA_DIR)

clean: clean-$(BUILD_DIR) clean-$(TEST_DATA_DIR)

test: clean
> @if [ ! -d $(TEST_DATA_DIR) ]; then
>   @mkdir -p $(TEST_TREE)
>   @touch $(TEST_FILES)
>   @$(FIRST_TEST_MSG)
>   @$(SECOND_TEST_MSG)
> @fi
>
> @go test -v ./...

.PHONY: build test clean
