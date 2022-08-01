#!/usr/bin/env bash

set -xv
set -euo pipefail

echo "compare:"
./compare-examples.sh

echo "explain:"
./explain-examples.sh

echo "show resources:"
./show-resources-examples.sh
