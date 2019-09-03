[![Build Status](https://travis-ci.org/redhat-nfvpe/service-assurance-poc.svg?branch=master)](https://travis-ci.org/redhat-nfvpe/helm2go-operator-sdk) [![Go Report Card](https://goreportcard.com/badge/github.com/redhat-nfvpe/helm2go-operator-sdk)](https://goreportcard.com/report/github.com/redhat-nfvpe/helm2go-operator-sdk)



## Overview
---
Helm2Go-Operator-SDK (H2Go) is a tool which creates the scaffold for  Go Operators corresponding to Helm Charts in a reproducible and scalable way. Read more about the design in the [design doc](docs/Design.md).

[Helm](https://github.com/helm/helm) is a tool used for managing Kubernetes charts. Charts are packages of pre-configured Kubernetes resources. Helm allows for versioning and distribution of native Kubernetes applications.

Go Operators are native Kubernetes applications used to deploy, upgrade, and manage other Kubernetes applications.

## Prerequisites
---
* [git](https://git-scm.com/downloads)
* [go](https://golang.org/dl/) version v1.12+
* [operator-sdk](https://github.com/operator-framework/operator-sdk) version v0.8+

## Compile and install from master

```sh
$ go get -d github.com/redhat-nfvpe/helm2go-operator-sdk # This will download the git repository and not install it
$ cd $GOPATH/src/github.com/redhat-nfvpe/helm2go-operator-sdk
$ git checkout master
$ make tidy
$ make dependency  # Installs operator sdk 
$ make install
```



## Quick Start
---
In the following example, we will create an nginx-operator using the existing [Bitnami Nginx](https://github.com/bitnami/charts/tree/master/bitnami/nginx) Helm Chart. 

### Create, Build and Deploy an *nginx-operator* from Local Chart
```
# Create an nginx-operator that defines the Ngnix CR
# Begin scaffolding process
$ helm2go-operator-sdk new nginx-operator --helm-chart=/path/to/nginx --api-version=web.example.com/v1alpha1 --kind=Ngnix
# Enter operator directory
$ cd nginx-operator

# Build the operator
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

*Note:* H2Go currently only supports non-deprecated resources. Thus, API Versions such as `apiVersion: extensions/v1beta1` should be updated to `apiVersion: apps/v1` or similar.
