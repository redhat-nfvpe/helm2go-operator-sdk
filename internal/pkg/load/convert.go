package load

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/validatemap"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes/scheme"
	//"k8s.io/api/extensions/v1beta1"
)

// YAMLUnmarshalResources converts a directory of injected YAML files to Kubernetes resources; resouresPath is assumed to be an absolute path
// if conversion of any resource fails method will panic; adds resources to resource cache
func YAMLUnmarshalResources(resourcesPath string, validMap *validatemap.ValidateMap, resourceCache *resourcecache.ResourceCache) (*resourcecache.ResourceCache, error) {

	// TODO convert all the files in a directory
	files, err := ioutil.ReadDir(resourcesPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// if the file is not contained in the map, continue without it
		if containValidMap(f.Name(), validMap) {
			logContinue()
			continue
		}

		file := filepath.Join(resourcesPath, f.Name())
		// obtain the converted result resourceKind
		rconfig, err := yamlUnmarshalSingleResource(file)
		if err != nil {
			if err.Error() == "empty" {
				continue
			} else if err.Error() == "not yaml" {
				continue
			} else if err.Error() == "unknown type" {
				return nil, fmt.Errorf("resource: %v is of unknown type", rconfig.resource)
			} else if err.Error() == "deprecated" {
				log.Printf("%s contains deprecated api version; exiting", f.Name())
				os.Exit(1)
			} else if err.Error() == "unsupported" {
				log.Printf("%s is an unsupported resource type; exiting", f.Name())
				os.Exit(1)
			} else {
				return nil, fmt.Errorf("unexpected error: %v", err)
			}
		}
		// add resource to cache
		// initializes the kind type
		resourceCache.SetKindType(rconfig.kindType)
		// set the correct cache information
		resourceCache.SetResourceForKindType(rconfig.kindType, rconfig.packageType)
		rs := resourceCache.GetResourceForKindType(rconfig.kindType)
		rs.SetResourceFunctions(string(rconfig.kindType), rconfig.resource)
	}
	return resourceCache, nil
}

// yamlUnmarshalSingleResource converts injected YAML files to a Kubernetes resource; rp is assumed to be an absolute path
func yamlUnmarshalSingleResource(rp string) (resourceConfig, error) {

	if filepath.Ext(rp) != ".yaml" && filepath.Ext(rp) != ".yml" {
		return resourceConfig{}, fmt.Errorf("not yaml")
	}

	// read and decode file
	fileBytes, err := ioutil.ReadFile(rp)
	if err != nil {
		return resourceConfig{}, fmt.Errorf("error reading file bytes from file: %s", rp)
	}
	return yamlUnmarshalSingleResourceFromBytes(fileBytes)
}

func yamlUnmarshalSingleResourceFromBytes(fileBytes []byte) (resourceConfig, error) {

	// // instantiate decoder for decoding purposes
	// baseScheme := scheme.Scheme
	// baseScheme.AddKnownTypes(networkingv1beta1.SchemeGroupVersion, &networkingv1beta1.Ingress{})
	// // if baseScheme.IsVersionRegistered(networkingv1beta1.SchemeGroupVersion) {
	// // 	return resourceConfig{}, fmt.Errorf("was not registered")
	// // }
	decode := scheme.Codecs.UniversalDeserializer().Decode
	if isEmptyFile(fileBytes) {
		return resourceConfig{}, fmt.Errorf("empty")
	}

	obj, _, err := decode([]byte(fileBytes), nil, nil)
	if err != nil {
		return resourceConfig{}, fmt.Errorf("error decoding bytes: %v", err)
	}

	// verify that the decoded resource kind is supported
	var kt resourcecache.KindType
	var pt resourcecache.PackageType
	 

	switch t := obj.(type) {
	case *corev1.Pod:
		kt = resourcecache.KindTypePod
		pt = resourcecache.PackageTypePods
	case *rbacv1.Role:
		kt = resourcecache.KindTypeRole
		pt = resourcecache.PackageTypeRoles
	case *rbacv1.RoleBinding:
		kt = resourcecache.KindTypeRoleBinding
		pt = resourcecache.PackageTypeRoleBindings
	case *rbacv1.ClusterRole:
		kt = resourcecache.KindTypeClusterRole
		pt = resourcecache.PackageTypeClusterRoles
	case *rbacv1.ClusterRoleBinding:
		kt = resourcecache.KindTypeClusterRoleBinding
		pt = resourcecache.PackageTypeClusterRoleBindings
	case *v1.Deployment:
		kt = resourcecache.KindTypeDeployment
		pt = resourcecache.PackageTypeDeployments
	case *corev1.Service:
		kt = resourcecache.KindTypeService
		pt = resourcecache.PackageTypeServices
	case *corev1.ServiceAccount:
		kt = resourcecache.KindTypeServiceAccount
		pt = resourcecache.PackageTypeServiceAccounts
	case *corev1.ConfigMap:
		kt = resourcecache.KindTypeConfigMap
		pt = resourcecache.PackageTypeConfigMaps	
	case *corev1.Secret:
		kt = resourcecache.KindTypeConfigMap
		pt = resourcecache.PackageTypeConfigMaps
	default:
		

		k := t.GetObjectKind().GroupVersionKind().Kind
		v := t.GetObjectKind().GroupVersionKind().Version

		if !AcceptedK8sTypes.MatchString(k) {
			return resourceConfig{
				obj,
				"",
				"",
			}, fmt.Errorf("unsupported")
		/*} else if v != "v1" {
			return resourceConfig{
				obj,
				"",
				"",
			}, fmt.Errorf("deprecated")*/
		} else {
			return resourceConfig{
				obj,
				"",
				"",
			}, fmt.Errorf("unknown type")
		}
	}
	rk := resourceConfig{
		obj,
		kt,
		pt,
	}

	return rk, nil
}
