#!/usr/bin/env bash

set -xv
set -euo pipefail


out_file=./resources-by-api-version-diff.md
KUBE_VERSIONS="1.20.15,1.21.14,1.22.17,1.23.17,1.24.12,1.25.8,1.26.3,1.27.0-rc.1"

echo '```' > $out_file
go run ../cmd/schema/main.go resources \
  --diff \
  --kube-version="${KUBE_VERSIONS}" \
  --group-by api-version >> $out_file
echo '```' >> $out_file
