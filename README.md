[![Build Status](https://travis-ci.org/redhat-nfvpe/service-assurance-poc.svg?branch=master)](https://travis-ci.org/redhat-nfvpe/helm2go-operator-sdk) [![Go Report Card](https://goreportcard.com/badge/github.com/redhat-nfvpe/helm2go-operator-sdk)](https://goreportcard.com/report/github.com/redhat-nfvpe/helm2go-operator-sdk)


![alt text](docs/design.png)

## Design

### Render
Render handles the primary steps in the Helm to Go Kubernetes pathway. The main responsibility of this package is to render valid Helm charts. Additionally, the package can write the injected files to a specified temp directory.

### Convert
Convert handles the secondary steps in the Helm to Go Kubernetes pathway. The main responsibility of this package is the unmarshal the rendered YAML files and produce raw Kubernetes resources.

The file `pkg/convert/convert.go` contains the main logic to accomplish the conversion itself. `YAMLUnmarshalResources` receives an absolute path to a directory and simply unmarshals all resource files one at a time.


## Flags
The existing flags are as listed:
```
--api-version string            Kubernetes apiVersion and has a format of $GROUP_NAME/$VERSION (e.g app.example.com/v1alpha1)
--cluster-scoped                Operator cluster scoped or not
--helm-chart string             Initialize helm operator with existing helm chart (<URL>, <repo>/<name>, or local path)
--helm-chart-ca-file string     CA File For External Repo (Optional)
--helm-chart-cert-file string   Cert File For External Repo (Optional)
--helm-chart-key-file string    Key File For External Repo (Optional)
--helm-chart-version string     Specific version of the helm chart (default is latest version)
--help                          help for convert
--kind string                   Kubernetes CustomResourceDefintion kind. (e.g AppService)
--password string               Password for chart repo (Optional)
--username string               Username for chart repo (Optional)
```

## How To Use

To create an operator from an existing *local* helm chart:
```
go run main.go convert <OperatorName> --helm-chart=/path/to/chart --kind=Kind --api-version=apps.example.com/v1alpha1
```

To create an operator from an *external* helm chart:
```
go run main.go convert <OperatorName> --helm-chart-repo=https://charts.example.io/ --helm-chart=example-chart --kind=Kind --api-version=apps.example.com/v1alpha1
```

## Common Problems
If you are experiencing build errors: `go: error loading module requirements`, execute the following command `export GO111MODULES=off` within the operator folder.