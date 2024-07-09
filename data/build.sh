#!/usr/bin/env bash

set -xv
set -euo pipefail


out_file=./resources-by-api-version-diff.md
KUBE_VERSIONS="1.24.17,1.25.16,1.26.15,1.27.15,1.28.11,1.29.6,1.30.2,1.31.0-alpha.3"

echo '```' > $out_file
go run ../cmd/schema/main.go resources \
  --diff \
  --kube-version="${KUBE_VERSIONS}" \
  --group-by api-version >> $out_file
echo '```' >> $out_file
