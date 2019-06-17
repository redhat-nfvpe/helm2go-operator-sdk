package resourcecache

import (
	"errors"
	"fmt"
	"strings"
)

//KindType ...
type KindType string

//PackageType ...
type PackageType string

//const for packagename types
const (
	PackageTypeRoles       PackageType = "roles"
	PackageTypeDeployments             = "deployments"
	PackageTypeContainers              = "containers"
	PackageTypePods                    = "pods"
	PackageTypeSecrets                 = "secrets"
	PackageTypeDaemonSet               = "daemonset"
	PackageTypeVolumes                 = "volumes"
	PackageTypeConfigMaps              = "configmaps"
)

//const for KindType
const (
	KindTypeConfigMap  KindType = "ConfigMap"
	KindTypeDeployment          = "Deployment"
	KindTypeRole                = "Role"
	KindTypeSecret              = "Secret"
	KindTypeVolume              = "Volume"
	KindTypeDaemonSet           = "DaemonSet"
	KindTypePod                 = "Pod"
	KindTypeContainer           = "Container"
)

/*
ResourceCache -->cache["ROLE"]->Resourc{package:roles,filename:deployment.go,functions:[]Functions}
*/
//ResourceCache ...
type ResourceCache struct {
	cache map[KindType]*Resource
}

//Resource ...
type Resource struct {
	PackageName PackageType
	FileName    string
	Functions   []ResourceFunction
}

//ResourceFunction  ...
type ResourceFunction struct {
	FunctionName string
	Data         string
}

//NewResourceCache  ...
func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		cache: make(map[KindType]*Resource),
	}
}

//Size  cache items
func (r *ResourceCache) Size() int {
	return len(r.cache)
}

//GetFunctionsByKind ....
func (r *ResourceCache) GetFunctionsByKind(kind KindType) []ResourceFunction {
	return r.cache[kind].Functions
}

//GetKindType ...
func (r *ResourceCache) GetKindType(kind KindType) (*Resource, error) {
	if value, ok := r.cache[kind]; ok {
		return value, nil
	}
	return nil, errors.New("Key not found ")

}

//SetKindType ...
func (r *ResourceCache) SetKindType(kind KindType) {
	if _, ok := r.cache[kind]; !ok {
		r.cache[kind] = &Resource{}
	}
}

//SetResourceForKindType ...
func (r *ResourceCache) SetResourceForKindType(kind KindType, packageType PackageType) {
	r.SetKindType(kind)
	filename := fmt.Sprintf("%s.go", strings.ToLower(string(kind)))
	if r.cache[kind] == nil {
		r.cache[kind] = &Resource{PackageName: packageType, FileName: filename}
	} else {
		r.cache[kind].FileName = filename
		r.cache[kind].PackageName = packageType
	}
}

//GetResourceForKindType ...
func (r *ResourceCache) GetResourceForKindType(kind KindType) *Resource {
	return r.cache[kind]
}

//GetResourceFunctions ...
func (rs *Resource) GetResourceFunctions() []ResourceFunction {
	return rs.Functions
}

//SetResourceFunctions ...
func (rs *Resource) SetResourceFunctions(name string, data string) {
	f := ResourceFunction{FunctionName: name, Data: data}
	rs.Functions = append(rs.Functions, f)
}

//Size  cache items
func (rs *Resource) Size() int {
	return len(rs.Functions)
}
