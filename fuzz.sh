#!/bin/bash

set -euxo pipefail

find . -name "*_test.go" \
    | xargs grep "func Fuzz" \
    | sed -E -e "s/^\.\/(([a-zA-Z0-9]+\/)+).+?(Fuzz[a-zA-Z0-9]+).+?$/\1 \3/g" \
    | while read package_name fuzz_test_name; do
        echo "$package_name $fuzz_test_name"
        "$GOROOT/bin/go" test --fuzz "$fuzz_test_name" --fuzztime=15s "github.com/bradenaw/juniper/$package_name"
    done
