package load

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

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

	for idx, f := range files {
		// obtain the converted result resourceKind
		fmt.Println(idx)
		rconfig, err := yamlUnmarshalSingleResource(filepath.Join(rp, f.Name()))
		if err != nil {
			if err.Error() == "Empty" || err.Error() == "Not YAML" {
				continue
			}
			return nil, err
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

func isEmptyFile(f []byte) bool {
	if len(f) == 0 {
		return true
	}
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	fmt.Println(string(f))
	if !isAlpha(string(f)) {
		return true
	}
	return false
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
		fmt.Println("We Empty")
		return resourceConfig{}, fmt.Errorf("Empty")
	}

	obj, resourceDefinition, err := decode([]byte(fileBytes), nil, nil)
	if err != nil {
		return resourceConfig{}, err
	}
	// verify that the decoded resource kind is supported
	k := resourceDefinition.Kind

	// TODO implement KubernetesResource
	var resource interface{}
	var kt resourcecache.KindType
	var pt resourcecache.PackageType

	switch k {
	case "Pod":
		resource = obj.(*corev1.Pod)
		kt = resourcecache.KindTypePod
		pt = resourcecache.PackageTypePods
	case "Role":
		resource = obj.(*rbacv1.Role)
		kt = resourcecache.KindTypeRole
		pt = resourcecache.PackageTypeRoles
	case "RoleBinding":
		resource = obj.(*rbacv1.RoleBinding)
		kt = resourcecache.KindTypeRoleBinding
		pt = resourcecache.PackageTypeRoleBindings
	case "Service":
		resource = obj.(*corev1.Service)
		kt = resourcecache.KindTypeService
		pt = resourcecache.PackageTypeServices
	case "Deployment":
		resource = obj.(*v1.Deployment)
		kt = resourcecache.KindTypeDeployment
		pt = resourcecache.PackageTypeDeployments
	default:
		log.Fatalf("Resource Type %v is Unsupported", k)
		return resourceConfig{}, err
	}
	rk := resourceConfig{
		resource,
		kt,
		pt,
	}
	return rk, nil

}
