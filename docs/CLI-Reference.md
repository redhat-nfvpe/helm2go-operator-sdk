# CLI Guide

**_Note:_** Binary has not yet been build.

```terminal
Usage:
    go run main.go [command]
```

## convert
---
Scaffolds a Go Operator for a corresponding Helm Chart.

### Args
* `operator-name` - name of the new operator

### Flags
* `--helm-chart` - Name of the helm chart. If using an external repo, specify the name within the repo i.e. `nginx`. If using a local chart provide `path/to/chart`
    * `--helm-chart-repo` - Specify external chart repo if necessary.
    * `--helm-chart-version` - Specify external chart version if necessary.
    * `--username` - Specify external repo username if necessary.
    * `--password` - Specify external repo password if necessary.
    * `--helm-chart-cert-file` - Specify Cert File for external repo if necessary.
    * `--helm-chart-key-file` - Specify Key File for external repo if necessary.
    * `--helm-chart-ca-file` - Specify CA File for external repo if necessary.
* `--api-version` - Kubernetes API Version and has a format of `<groupName>/<version>` i.e. `app.example.com/v1alpha1`
* `--kind` - Kubernetes Custom Resource Definition kind.
* `--cluster-scoped` - Operator cluster scoped or not.

### Example
```
$ go run main.go convert nginx-operator --helm-chart=path/to/chart --api-version=web.example.com/v1alpha1 --kind=Nginx
```