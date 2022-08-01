# kubectl-schema

## Overview

This program uses kubernetes openapi swagger specs to help you understand what Resources and ApiVersions
are supported by kubernetes versions.

This is a multidimensional calculation, across:

 - Resources (example: Pod, Ingress)
 - Api Versions (example: v1, batch/v1, apiextensions.k8s.io/v1beta1)
 - Kubernetes versions
 - Resource schema

## Installation

### Golang binary

TODO

### Kubectl plugin

TODO

## Examples

### Resources

#### What api-versions are resources in?

```
kubectl schema resources \
  --resource=Ingress,CronJob,CustomResourceDefinition \
  --kube-version=1.18.20,1.20.15,1.22.12,1.24.0,1.25.0-alpha.3

+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
|         RESOURCE         |           1.18.20            |           1.20.15            |         1.22.12         |         1.24.0          |     1.25.0-ALPHA 3      |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| CronJob                  | batch.v1beta1                | batch.v1beta1                | batch.v1                | batch.v1                | batch.v1                |
|                          | batch.v2alpha1               | batch.v2alpha1               | batch.v1beta1           | batch.v1beta1           |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| CustomResourceDefinition | apiextensions.k8s.io.v1      | apiextensions.k8s.io.v1      | apiextensions.k8s.io.v1 | apiextensions.k8s.io.v1 | apiextensions.k8s.io.v1 |
|                          | apiextensions.k8s.io.v1beta1 | apiextensions.k8s.io.v1beta1 |                         |                         |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| Ingress                  | extensions.v1beta1           | extensions.v1beta1           | networking.k8s.io.v1    | networking.k8s.io.v1    | networking.k8s.io.v1    |
|                          | networking.k8s.io.v1beta1    | networking.k8s.io.v1         |                         |                         |                         |
|                          |                              | networking.k8s.io.v1beta1    |                         |                         |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
```

#### What resources are in api-versions?

```
kubectl schema resources \
  --resource=Ingress,CronJob,CustomResourceDefinition \
  --kube-version=1.18.20,1.20.15,1.22.12,1.24.0,1.25.0-alpha.3 \
  --group-by=api-version

+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
|         API VERSION          |         1.18.20          |         1.20.15          |         1.22.12          |          1.24.0          |      1.25.0-ALPHA 3      |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| apiextensions.k8s.io.v1      | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition |
+------------------------------+                          +                          +--------------------------+--------------------------+--------------------------+
| apiextensions.k8s.io.v1beta1 |                          |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| batch.v1                     |                          |                          | CronJob                  | CronJob                  | CronJob                  |
+------------------------------+--------------------------+--------------------------+                          +                          +--------------------------+
| batch.v1beta1                | CronJob                  | CronJob                  |                          |                          |                          |
+------------------------------+                          +                          +--------------------------+--------------------------+--------------------------+
| batch.v2alpha1               |                          |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| extensions.v1beta1           | Ingress                  | Ingress                  |                          |                          |                          |
+------------------------------+--------------------------+                          +--------------------------+--------------------------+--------------------------+
| networking.k8s.io.v1         |                          |                          | Ingress                  | Ingress                  | Ingress                  |
+------------------------------+--------------------------+                          +--------------------------+--------------------------+--------------------------+
| networking.k8s.io.v1beta1    | Ingress                  |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
```

#### Show the change in api-versions for given resources

```
kubectl schema resources \
  --resource=Ingress,CronJob,CustomResourceDefinition \
  --kube-version=1.18.20,1.20.15,1.22.12,1.24.0,1.25.0-alpha.3 \
  --diff

+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
|         RESOURCE         |           1.18.20            |        1.20.15         |            1.22.12             | 1.24.0 | 1.25.0-ALPHA 3  |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| CronJob                  | batch.v1beta1                |                        | remove:                        |        | remove:         |
|                          | batch.v2alpha1               |                        |   batch.v2alpha1               |        |   batch.v1beta1 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| CustomResourceDefinition | apiextensions.k8s.io.v1      |                        | remove:                        |        |                 |
|                          | apiextensions.k8s.io.v1beta1 |                        |   apiextensions.k8s.io.v1beta1 |        |                 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| Ingress                  | extensions.v1beta1           | add:                   | remove:                        |        |                 |
|                          | networking.k8s.io.v1beta1    |   networking.k8s.io.v1 |   extensions.v1beta1           |        |                 |
|                          |                              |                        |   networking.k8s.io.v1beta1    |        |                 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
```

#### Show the change in resources for given api-versions

```
kubectl schema resources \
  --resource=Ingress,CronJob,CustomResourceDefinition \
  --kube-version=1.18.20,1.20.15,1.22.12,1.24.0,1.25.0-alpha.3 \
  --group-by=api-version \
  --diff

+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
|         API VERSION          |         1.18.20          |  1.20.15  |          1.22.12           | 1.24.0 | 1.25.0-ALPHA 3 |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| apiextensions.k8s.io.v1      | CustomResourceDefinition |           |                            |        |                |
+------------------------------+                          +-----------+----------------------------+--------+----------------+
| apiextensions.k8s.io.v1beta1 |                          |           | remove:                    |        |                |
|                              |                          |           |   CustomResourceDefinition |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| batch.v1                     |                          |           | add:                       |        |                |
|                              |                          |           |   CronJob                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| batch.v1beta1                | CronJob                  |           |                            |        | remove:        |
|                              |                          |           |                            |        |   CronJob      |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+                          +-----------+----------------------------+--------+----------------+
| batch.v2alpha1               |                          |           | remove:                    |        |                |
|                              |                          |           |   CronJob                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| extensions.v1beta1           | Ingress                  |           | remove:                    |        |                |
|                              |                          |           |   Ingress                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| networking.k8s.io.v1         |                          | add:      |                            |        |                |
|                              |                          |   Ingress |                            |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| networking.k8s.io.v1beta1    | Ingress                  |           | remove:                    |        |                |
|                              |                          |           |   Ingress                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
```

### Explain

### Compare

## Dev

### How to release a new binary

See `goreleaser`'s requirements [here](https://goreleaser.com/environment/).

Get a [GitHub Personal Access Token](https://github.com/settings/tokens/new) and add the `repo` scope.
Set `GITHUB_TOKEN` to this value:

```bash
export GITHUB_TOKEN=...
```

[See here for more information on github tokens](https://help.github.com/articles/creating-an-access-token-for-command-line-use/).

Choose a tag/release name, create and push a tag:

```bash
TAG=v0.0.1

git tag $TAG
git push origin $TAG
```

Cut a release:

```bash
goreleaser release --rm-dist
```

Make a test release:

```bash
goreleaser release --snapshot --rm-dist
```
