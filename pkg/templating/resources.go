package templating

import "reflect"

// Imports is a map of the imports needed for respective resource types
var Imports = map[string][]string{

	// TODO continue to fill in the other imports
	// only focusing on deployment right now

	"ConfigMap": []string{
		`corev1 "k8s.io/api/core/v1`,
		`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
	},
	"Deployment": []string{
		`appsv1 "k8s.io/api/apps/v1"`,
		`corev1 "k8s.io/api/core/v1"`,
		`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
	},
	"Secret":    []string{},
	"Volume":    []string{},
	"DaemonSet": []string{},
	"Pod":       []string{},
	"Container": []string{
		`corev1 "k8s.io/api/core/v1"`,
	},
	"Service": []string{},
}

func getTemplateImports(resource interface{}) []string {
	// get the kind from the resource
	inferredResourceTypeName := reflect.TypeOf(resource).Name()
	// access the Imports map with the corresponding key
	resourceTypeImports := Imports[inferredResourceTypeName]
	// return the Imports value
	return resourceTypeImports
}

func Test(resource interface{}) []string {
	return getTemplateImports(resource)
}
