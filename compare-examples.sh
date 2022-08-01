#!/usr/bin/env bash

set -xv
set -euo pipefail


go run cmd/schema/main.go compare \
  --kube-version 1.18.19,1.24.3 \
  --resource CustomResourceDefinition,CronJob #> compare-crd.txt

go run cmd/schema/main.go compare \
  --kube-version 1.18.0,1.24.2 \
  --resource NetworkPolicy,Ingress
