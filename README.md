[![Build Status](https://travis-ci.org/redhat-nfvpe/service-assurance-poc.svg?branch=master)](https://travis-ci.org/redhat-nfvpe/helm2go-operator-sdk) [![Go Report Card](https://goreportcard.com/badge/github.com/redhat-nfvpe/helm2go-operator-sdk)](https://goreportcard.com/report/github.com/redhat-nfvpe/helm2go-operator-sdk)



## Overview
---
This project is a small tool to produce Go Operators corresponding to Helm Charts in a reproducible and scalable way. Read more about the design in the [design doc](docs/Design.md).

[Helm](https://github.com/helm/helm) is a tool used for managing Kubernetes charts. Charts are packages of pre-configured Kubernetes resources. Helm allows for versioning and distribution of native Kubernetes applications.

Go Operators are native Kubernetes applications used to deploy, upgrade, and manage other Kubernetes applications.


## Workflow
---
The tool provides the following workflow to develop operators in Go from corresponding Helm Charts:

1. Identify Helm Chart
    * Tool supports local charts
    * Tool supports external charts i.e. those hosted on external repositories
2. Specify the neccessary resource, and the API is generated adding Custom Resource Definitions (CRDs)
3. *Supported* Kubernetes Resources Controllers are automatically generated
4. User must write the reconciling logic for the controller using the [Operator-SDK](https://github.com/operator-framework/operator-sdk) and [controller-runtime](https://godoc.org/sigs.k8s.io/controller-runtime) APIs. 
5. Use the [Operator-SDK](https://github.com/operator-framework/operator-sdk) CLI to build and generate the operator deployment manifests. 


## Prerequisites
---
* [git](https://git-scm.com/downloads)
* [go](https://golang.org/dl/) version v1.12+
* [operator-sdk](https://github.com/operator-framework/operator-sdk) version v0.8+
* [dep](https://golang.github.io/dep/docs/installation.html) version v0.5.0+


## Quick Start
---
In the following example, we will create an nginx-operator using the existing [Bitnami Nginx](https://github.com/bitnami/charts/tree/master/bitnami/nginx) Helm Chart. 

### Create, Build and Deploy an *nginx-operator* from Local Chart
```
# Create an nginx-operator that defines the Ngnix CR
$ export GO111MODULE=on
# Begin scaffolding process
$ helm2go-operator-sdk convert nginx-operator --helm-chart=/path/to/nginx --api-version=web.example.com/v1alpha1 --kind=Ngnix
# Enter operator directory
$ cd nginx-operator

# Build the operator
$ export GO111MODULE=off
$ operator-sdk build quay.io/example/image
$ docker push quay.io/example/image

# Deploy the Operator
# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
$ sed -i 's|REPLACE_IMAGE|quay.io/example/app-operator|g' deploy/operator.yaml
# On OSX use:
$ sed -i "" 's|REPLACE_IMAGE|quay.io/example/app-operator|g' deploy/operator.yaml

# Setup Service Account
$ kubectl create -f deploy/service_account.yaml
# Setup RBAC
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
# Setup the CRD
$ kubectl create -f deploy/crds/web_v1alpha1_nginx_crd.yaml
# Deploy the app-operator
$ kubectl create -f deploy/operator.yaml

# Create an AppService CR
# The default controller will watch for AppService objects and create a pod for each CR
$ kubectl create -f deploy/crds/app_v1alpha1_nginx_cr.yaml

# Verify that a pod is created
$ kubectl get pod -l app=example-nginx
NAME                     READY     STATUS    RESTARTS   AGE
example-nginx-pod   1/1       Running   0          1m
```

## Supported Kubernetes Resources
---
The tool currently supports a very limited selection of Kubernetes resources, as listed here:
```
Role
ClusterRole
RoleBinding
ClusterRoleBinding
ServiceAccount
Service
Deployment
```
If attempting to parse a Kubernetes resource other than the ones listed above the tool will prompt the user to either `continue` code generation without the unsupported resources or `stop` the code generation all together.


## Common Problems
---
If you are experiencing build errors: `go: error loading module requirements`, execute the following command `export GO111MODULES=off` within the operator folder.