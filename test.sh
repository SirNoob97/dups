#!/bin/env bash

set -eu

test_dir="./testdata"

function setup_test_dir {
    mkdir -p $test_dir/{1..5}/{a..e}
    touch $test_dir/{test_{1..5},empty_file,{1..5}/{test_{1..5},{a..e}/test_{a..e}}}
    echo "this is a test file" | tee $test_dir/{test_{1..3},{1..5}/{test_{1..3},{a..e}/test_{a..c}}} 1>/dev/null
    echo "this is a another test file" | tee $test_dir/{test_{4,5},{1..5}/{test_{4,5},{a..e}/test_{d,e}}} 1>/dev/null
}

if [ ! -d $test_dir ]; then
    setup_test_dir
fi

go test -v ./...

exit 0
