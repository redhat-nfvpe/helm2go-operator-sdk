package templating

import (
	"reflect"
)

// Imports is a map of the imports needed for respective resource types
var Imports = map[string][]string{

	// TODO continue to fill in the other imports
	// only focusing on deployment right now

	"*v1.ConfigMap": []string{
		`corev1 "k8s.io/api/core/v1`,
		`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
	},
	"*v1.Deployment": []string{
		`appsv1 "k8s.io/api/apps/v1"`,
		`corev1 "k8s.io/api/core/v1"`,
		`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
	},
}

// PackageNames is the cannonical map for creating package names
var PackageNames = map[string]string{
	"*v1.Deployment": "deployments",
}

// ResourceTitles is the cannonical map for pretty resource names
var ResourceTitles = map[string]string{
	"*v1.Deployment": "Deployment",
}

func getTemplateImports(resource *interface{}) []string {
	inferredResourceTypeName := reflect.TypeOf(*resource).String()
	resourceTypeImports := Imports[inferredResourceTypeName]
	return resourceTypeImports
}

func getTemplatePackageName(resource *interface{}) string {
	inferredResourceTypeName := reflect.TypeOf(*resource).String()
	resourcePackageName := PackageNames[inferredResourceTypeName]
	return resourcePackageName
}

func getTemplateResourceTitle(resource *interface{}) string {
	inferredResourceTypeName := reflect.TypeOf(*resource).String()
	resourceTitle := ResourceTitles[inferredResourceTypeName]
	return resourceTitle
}
