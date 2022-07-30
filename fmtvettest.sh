#!/usr/bin/env bash

set -xv
set -euo pipefail

make fmt

make vet

make test
