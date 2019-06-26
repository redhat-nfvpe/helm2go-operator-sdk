package resourcecache

import (
	"errors"
	"fmt"
	"strings"
)

//KindType ...
type KindType string

func (k *KindType) String() string {
	return string(*k)
}

//PackageType ...
type PackageType string

func (p *PackageType) String() string {
	return string(*p)
}

//const for packagename types
const (
	PackageTypeRoles               PackageType = "roles"
	PackageTypeClusterRoles                    = "clusterroles"
	PackageTypeDeployments                     = "deployments"
	PackageTypeContainers                      = "containers"
	PackageTypePods                            = "pods"
	PackageTypeSecrets                         = "secrets"
	PackageTypeDaemonSet                       = "daemonset"
	PackageTypeVolumes                         = "volumes"
	PackageTypeConfigMaps                      = "configmaps"
	PackageTypeServices                        = "services"
	PackageTypeRoleBindings                    = "rolebindings"
	PackageTypeClusterRoleBindings             = "clusterrolebindings"
	PackageTypeServiceAccounts                 = "serviceaccounts"
)

//const for KindType
const (
	KindTypeConfigMap          KindType = "ConfigMap"
	KindTypeDeployment                  = "Deployment"
	KindTypeRole                        = "Role"
	KindTypeClusterRole                 = "ClusterRole"
	KindTypeSecret                      = "Secret"
	KindTypeVolume                      = "Volume"
	KindTypeDaemonSet                   = "DaemonSet"
	KindTypePod                         = "Pod"
	KindTypeContainer                   = "Container"
	KindTypeService                     = "Service"
	KindTypeRoleBinding                 = "RoleBinding"
	KindTypeClusterRoleBinding          = "ClusterRoleBinding"
	KindTypeServiceAccount              = "ServiceAccount"
)

/*
ResourceCache -->cache["ROLE"]->Resourc{package:roles,filename:roles.go,functions:[]Functions}
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
	Data         interface{}
}

//NewResourceCache  ...
func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		cache: make(map[KindType]*Resource),
	}
}

// GetCache returns the main cache object
func (r *ResourceCache) GetCache() *map[KindType]*Resource {
	return &r.cache
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
	return nil, errors.New("key not found")

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

//PrepareCacheForFile formats the cache for the file write operation
func (r *ResourceCache) PrepareCacheForFile() map[string]*Resource {
	c := *r.GetCache()
	var newKey string
	var newMap map[string]*Resource
	newMap = make(map[string]*Resource)
	for rtype := range c {
		// reset the key to be a filename
		newKey = nameToFileName(rtype.String(), FileExtensionGo)
		newMap[newKey] = c[rtype]
	}
	return newMap
}

//GetKindTypes returns a list of kind types in use by the cache
func (r *ResourceCache) GetKindTypes() []KindType {
	kts := make([]KindType, len(r.cache))
	i := 0
	for kt := range r.cache {
		kts[i] = kt
		i++
	}
	return kts
}

//GetResourceFunctions ...
func (rs *Resource) GetResourceFunctions() []ResourceFunction {
	return rs.Functions
}

//SetResourceFunctions ...
func (rs *Resource) SetResourceFunctions(functionname string, data interface{}) {
	f := ResourceFunction{FunctionName: functionname, Data: data}
	rs.Functions = append(rs.Functions, f)
}

//Size  cache items
func (rs *Resource) Size() int {
	return len(rs.Functions)
}
