package convert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1" //"k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// JSONUnmarshal receives a resource json file(deployment for right now) and returns a Go Kubernetes resource (only deployment right now)
func JSONUnmarshal(path string) (*v1.Deployment, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dep v1.Deployment
	err = json.Unmarshal(bytes, &dep)
	if err != nil {
		return nil, err
	}
	return &dep, nil
}

// DirectoryInjectedYAMLToJSON ...
func DirectoryInjectedYAMLToJSON(resourcesPath string) {
	// read the directory of injected YAML files
	files, err := ioutil.ReadDir(resourcesPath)
	if err != nil {
		log.Fatal(err)
	}
	// instantiate decoder for use in deserializing YAML files
	decode := scheme.Codecs.UniversalDeserializer().Decode

	acceptedK8sTypes := regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount|Service|Deployment)`)

	// iterate over the YAML files
	for idx, f := range files {
		fileBytes, err := ioutil.ReadFile(filepath.Join("./test/resources/", f.Name()))
		obj, groupVersionKind, err := decode([]byte(fileBytes), nil, nil)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
		}
		// TODO the following file path is hard-coded, need to change it to be robust
		if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			fmt.Println("Error while decoding, its not in acceptable list")
		} else {

			var resource interface{}

			switch groupVersionKind.Kind {
			//case *corev1.Pod:
			case "Pod":
				resource = obj.(*corev1.Pod)
			case "Role":
				resource = obj.(*rbacv1.Role)
			case "RoleBinding":
				resource = obj.(*rbacv1.RoleBinding)
			case "ClusterRole":
				resource = obj.(*rbacv1.ClusterRole)
			case "ClusterRoleBinding":
				resource = obj.(*rbacv1.ClusterRoleBinding)
			case "ServiceAccount":
				resource = obj.(*corev1.ServiceAccount)
			case "Service":
				resource = obj.(*corev1.Service)
			case "Deployment":
				resource = obj.(*v1.Deployment)
			default:
				fmt.Printf("%+v\n", groupVersionKind.Kind)
				fmt.Println("Unknown Kind")
			}

			// need to pretty print the resource object for use in template
			f := fmt.Sprintf("./test/jsonOutputs/test_%s_%d.json", groupVersionKind.Kind, idx)
			file, _ := json.MarshalIndent(resource, "", " ")
			_ = ioutil.WriteFile(f, file, 0644)
		}
	}
}
