#!/bin/bash

set -euo pipefail

go_versions=(1.21.9 1.22.9 1.23.6 1.24.0)

latest="${go_versions[-1]}"
if ! go version | grep "go$latest" > /dev/null; then
    echo >&2 "go version expected $latest, got $(go version)"
    exit 1
fi

for go_version in ${go_versions[@]}; do
    if [[ ${go_version} == ${latest} ]]; then
        go version
        go test --race ./...
    else
        go install "golang.org/dl/go${go_version}@latest"
        go_bin="${HOME}/go/bin/go${go_version}"
        $go_bin download
        $go_bin version
        $go_bin test --race ./...
    fi
done
