#!/usr/bin/env bash

set -xv
set -euo pipefail


out_file=./resources-by-api-version-diff.md

echo '```' > $out_file
go run ../cmd/schema/main.go resources --diff --group-by api-version >> $out_file
echo '```' >> $out_file
