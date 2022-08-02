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

Download the latest binary for your platform from https://github.com/mattfenwick/kubectl-schema/releases .

### Kubectl plugin

After downloading a `kubectl-schema` binary, place it somewhere in your `$PATH`.
It will now be usable as `kubectl schema`.

See [the kubectl docs](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) for more information.

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

#### Use a path to focus results

```bash
kubectl schema explain \
  --kube-version 1.24.0 \
  --resource Ingress \
  --format table \
  --path spec.tls,status

1.24.0
io.k8s.api.networking.v1 Ingress:
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
|                       PATH                       |                                                 TYPE                                                 |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| spec.tls                                         | array                                                                                                |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| spec.tls.[]                                      | object                                                                                               |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| spec.tls.[].hosts                                | array                                                                                                |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| spec.tls.[].hosts.[]                             | string                                                                                               |
+--------------------------------------------------+                                                                                                      +
| spec.tls.[].secretName                           |                                                                                                      |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status                                           | object                                                                                               |
+--------------------------------------------------+                                                                                                      +
| status.loadBalancer                              |                                                                                                      |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress                      | array                                                                                                |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[]                   | object                                                                                               |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].hostname          | string                                                                                               |
+--------------------------------------------------+                                                                                                      +
| status.loadBalancer.ingress.[].ip                |                                                                                                      |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].ports             | array                                                                                                |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].ports.[]          | object                                                                                               |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].ports.[].error    | string                                                                                               |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].ports.[].port     | integer                                                                                              |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
| status.loadBalancer.ingress.[].ports.[].protocol | string                                                                                               |
+--------------------------------------------------+------------------------------------------------------------------------------------------------------+
```

#### Use a `--depth` to limit results

```bash
kubectl schema explain \
  --kube-version 1.24.0 \
  --resource Ingress \
  --format table \
  --depth 1              

1.24.0
io.k8s.api.networking.v1 Ingress:
+------------+------------------------------------------------------------------------------------------------------+
|    PATH    |                                                 TYPE                                                 |
+------------+------------------------------------------------------------------------------------------------------+
|            | object                                                                                               |
+------------+------------------------------------------------------------------------------------------------------+
| apiVersion | string                                                                                               |
+------------+                                                                                                      +
| kind       |                                                                                                      |
+------------+------------------------------------------------------------------------------------------------------+
| metadata   | object                                                                                               |
+------------+                                                                                                      +
| spec       |                                                                                                      |
+------------+                                                                                                      +
| status     |                                                                                                      |
+------------+------------------------------------------------------------------------------------------------------+
```

### Compare

Compare the schema for a type between multiple kubernetes versions.

```bash
kubectl schema compare \
  --kube-version 1.18.0,1.24.2 \
  --resource Ingress      

comparing Ingress: 1.18.0@io.k8s.api.extensions.v1beta1 vs. 1.24.2@io.k8s.api.networking.v1
  +                       metadata.managedFields.[].subresource
  -                       spec.backend
  -                       spec.rules.[].http.paths.[].backend.serviceName
  -                       spec.rules.[].http.paths.[].backend.servicePort
  +                       spec.rules.[].http.paths.[].backend.service
  +                       spec.defaultBackend
  +                       status.loadBalancer.ingress.[].ports

comparing Ingress: 1.18.0@io.k8s.api.networking.v1beta1 vs. 1.24.2@io.k8s.api.networking.v1
  +                       metadata.managedFields.[].subresource
  -                       spec.backend
  -                       spec.rules.[].http.paths.[].backend.serviceName
  -                       spec.rules.[].http.paths.[].backend.servicePort
  +                       spec.rules.[].http.paths.[].backend.service
  +                       spec.defaultBackend
  +                       status.loadBalancer.ingress.[].ports
```

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
