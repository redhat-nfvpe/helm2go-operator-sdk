package load

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes/scheme"
	//"k8s.io/api/extensions/v1beta1"
)

// YAMLUnmarshalResources converts a directory of injected YAML files to Kubernetes resources; resouresPath is assumed to be an absolute path
// if conversion of any resource fails method will panic; adds resources to resource cache
func YAMLUnmarshalResources(rp string, rc *resourcecache.ResourceCache) (*resourcecache.ResourceCache, error) {

	// TODO convert all the files in a directory
	files, err := ioutil.ReadDir(rp)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		// obtain the converted result resourceKind
		rconfig, err := yamlUnmarshalSingleResource(filepath.Join(rp, f.Name()))
		if err != nil {
			if err.Error() == "Empty" || err.Error() == "Not YAML" || err.Error() == "Unsupported" {
				continue
			}
		}
		// add resource to cache
		// initializes the kind type
		rc.SetKindType(rconfig.kt)
		// set the correct cache information
		rc.SetResourceForKindType(rconfig.kt, rconfig.pt)
		rs := rc.GetResourceForKindType(rconfig.kt)
		rs.SetResourceFunctions(string(rconfig.kt), rconfig.r)
	}
	return rc, nil
}

// yamlUnmarshalSingleResource converts injected YAML files to a Kubernetes resource; rp is assumed to be an absolute path
func yamlUnmarshalSingleResource(rp string) (resourceConfig, error) {

	if filepath.Ext(rp) != ".yaml" && filepath.Ext(rp) != ".yml" {
		return resourceConfig{}, fmt.Errorf("Not YAML")
	}

	// instantiate decoder for decoding purposes
	decode := scheme.Codecs.UniversalDeserializer().Decode
	// read and decode file
	fileBytes, err := ioutil.ReadFile(rp)

	if isEmptyFile(fileBytes) {
		return resourceConfig{}, fmt.Errorf("Empty")
	}

	obj, _, err := decode([]byte(fileBytes), nil, nil)
	if err != nil {
		return resourceConfig{}, err
	}

	// verify that the decoded resource kind is supported
	var kt resourcecache.KindType
	var pt resourcecache.PackageType

	switch obj.(type) {
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
	default:
		log.Printf("Resource Type %v is Unsupported! Update API Version to Continue.", reflect.TypeOf(obj))
		return resourceConfig{}, fmt.Errorf("Unsupported")
	}
	rk := resourceConfig{
		obj,
		kt,
		pt,
	}

	return rk, nil

}
