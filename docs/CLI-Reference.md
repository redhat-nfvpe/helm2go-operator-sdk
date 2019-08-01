# CLI Guide

```terminal
Usage:
    helm2go-operator-sdk [command] [arguments] [flags]
```

## new
---
Scaffolds a Go Operator for a corresponding Helm Chart.

### Args
* `operator-name` - name of the new operator

### Flags
* **Required:**` --helm-chart` - Name of the helm chart. If using an external repo, specify the name within the repo i.e. `nginx`. If using a local chart provide `path/to/chart`
* `--helm-chart-repo` - Specify external chart repo if necessary.
* `--helm-chart-version` - Specify external chart version if necessary.
* `--username` - Specify external repo username if necessary.
* `--password` - Specify external repo password if necessary.
* `--helm-chart-cert-file` - Specify Cert File for external repo if necessary.
* `--helm-chart-key-file` - Specify Key File for external repo if necessary.
* `--helm-chart-ca-file` - Specify CA File for external repo if necessary.
* **Required:** `--api-version` - Kubernetes API Version and has a format of `<groupName>/<version>` i.e. `app.example.com/v1alpha1`
* **Required:** `--kind` - Kubernetes Custom Resource Definition kind.
* `--cluster-scoped` - Operator cluster scoped or not.

### Example
```
$ helm2go-operator-sdk new nginx-operator --helm-chart=path/to/chart --api-version=web.example.com/v1alpha1 --kind=Nginx
```

The resulting structure will be:
```
<project-name>
|   build/
|   |   Dockerfile
|   |   bin/
|   |   _output/
|   cmd/.
|   |   manager/
|   |   |   main.go
|   deploy/
|   |   operator.yaml
|   |   role_binding.yaml
|   |   role.yaml
|   |   service_account.yaml
|   |   crds/
|   |   |   <api-version>_<kind>_cr.yaml
|   |   |   <api-version>_<kind>_crd.yaml
|   pkg/
|   |   apis/
|   |   |   <api-version>/
|   |   |   |   doc.go
|   |   |   |   memcached_types.go
|   |   |   |   ...
|   |   |   ...
|   |   controller/
|   |   |   <kind>/
|   |   |   <kind>_controller.go
|   |   |   ...
|   |   resources/
|   |   |   ...
|   ...
```

