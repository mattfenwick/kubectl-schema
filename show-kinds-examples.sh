#!/usr/bin/env bash

set -xv
set -euo pipefail


# kinds: explain
KUBE_VERSIONS=1.18.20,1.20.15,1.22.12,1.24.0,1.25.0-alpha.3
KUBE_RESOURCES=Ingress,CronJob,CustomResourceDefinition
go run cmd/schema/main.go kinds \
  --resource="$KUBE_RESOURCES" \
  --kube-version="$KUBE_VERSIONS"

go run cmd/schema/main.go kinds \
  --resource="$KUBE_RESOURCES" \
  --kube-version="$KUBE_VERSIONS" \
  --group-by=api-version

go run cmd/schema/main.go kinds \
  --resource="$KUBE_RESOURCES" \
  --kube-version="$KUBE_VERSIONS" \
  --diff

go run cmd/schema/main.go kinds \
  --resource="$KUBE_RESOURCES" \
  --kube-version="$KUBE_VERSIONS" \
  --group-by=api-version \
  --diff
