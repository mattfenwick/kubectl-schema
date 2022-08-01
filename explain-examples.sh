#!/usr/bin/env bash

set -xv
set -euo pipefail


go run cmd/schema/main.go explain \
  --kube-version 1.18.19 \
  --resource CustomResourceDefinition #> explain-crd.txt

go run cmd/schema/main.go explain \
  --kube-version 1.18.19 \
  --resource CronJob > cronjob-1-18.txt

go run cmd/schema/main.go explain \
  --kube-version 1.24.3 \
  --resource CronJob > cronjob-1-24.txt
