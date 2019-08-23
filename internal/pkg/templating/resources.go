package templating

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// TemplateResource is for organizational purposes
type TemplateResource struct {
	resourceType string
	kindType     resourcecache.KindType
	packageType  resourcecache.PackageType
}

// KindTypeString returns the kind type
func (t *TemplateResource) KindTypeString() string {
	return t.kindType.String()
}

// PackageTypeString returns the package type
func (t *TemplateResource) PackageTypeString() string {
	return t.packageType.String()
}

// // TemplateConfig is an interface for all template configuration objects
// type TemplateConfig interface {
// 	Execute()
// 	GetTemplate()
// }

var importStatements = []string{
	`appsv1 "k8s.io/api/apps/v1"`,
	`corev1 "k8s.io/api/core/v1"`,
	`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
}
var rbacImportStatements = []string{
	`corev1 "k8s.io/api/core/v1"`,
	`metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`,
	`rbacv1 "k8s.io/api/rbac/v1"`,
}

//ImportPackages is the cannonical map for getting the correct import package
var ImportPackages = map[string]string{
	"Deployment":         "appsv1",
	"ConfigMap":          "corev1",
	"Container":          "corev1",
	"Service":            "corev1",
	"Pod":                "corev1",
	"Secret":             "corev1",
	"Volume":             "corev1",
	"ServiceAccount":     "corev1",
	"Role":               "rbacv1",
	"ClusterRole":        "rbacv1",
	"RoleBinding":        "rbacv1",
	"ClusterRoleBinding": "rbacv1",
}

// PackageNames is the cannonical map for creating package names
var PackageNames = map[string]string{
	"*v1.Deployment":         "deployments",
	"*v1.Service":            "services",
	"*v1.ConfigMap":          "configmaps",
	"*v1.ServiceAccount":     "serviceaccount",
	"*v1.ClusterRole":        "clusterroles",
	"*v1.ClusterRoleBinding": "clusterrolebindings",
	"*v1.Role":               "roles",
	"*v1.RoleBinding":        "rolebindings",
}

// ResourceTitles is the cannonical map for pretty resource names
var ResourceTitles = map[string]string{
	"*v1.Deployment":         "Deployment",
	"*v1.Service":            "Service",
	"*v1.ConfigMap":          "ConfigMap",
	"*v1.ServiceAccount":     "ServiceAccount",
	"*v1.ClusterRole":        "ClusterRole",
	"*v1.ClusterRoleBinding": "ClusterRoleBinding",
	"*v1.Role":               "Role",
	"*v1.RoleBinding":        "RoleBinding",
}

var controllerKindImports = map[string]string{
	"k8s.io/api/core/v1":                                           "corev1",
	"k8s.io/api/apps/v1":                                           "appsv1",
	"k8s.io/apimachinery/pkg/api/errors":                           "",
	"k8s.io/apimachinery/pkg/apis/meta/v1":                         "metav1",
	"k8s.io/apimachinery/pkg/runtime":                              "",
	"k8s.io/apimachinery/pkg/types":                                "",
	"sigs.k8s.io/controller-runtime/pkg/client":                    "",
	"sigs.k8s.io/controller-runtime/pkg/controller":                "",
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil": "",
	"sigs.k8s.io/controller-runtime/pkg/handler":                   "",
	"sigs.k8s.io/controller-runtime/pkg/manager":                   "",
	"sigs.k8s.io/controller-runtime/pkg/reconcile":                 "",
	"sigs.k8s.io/controller-runtime/pkg/runtime/log":               "logf",
	"sigs.k8s.io/controller-runtime/pkg/source":                    "",
}

func getTemplateImports(resource *interface{}) []string {
	return importStatements
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

func getResourceImportPackage(resourceType string) string {
	importName := ImportPackages[resourceType]
	return importName
}

func kindToLower(kind string) string {
	return strings.ToLower(kind)
}

func getNamePlural(input string) string {
	if input[len(input)-1] == 's' {
		return input + "es"
	}
	return input + "s"
}

func kindToLowerCamel(kind string) string {
	return strcase.ToLowerCamel(kind)
}

func getOwnerAPIVersion(apiVersion, kind string) string {
	return kindToLowerCamel(kind) + filepath.Base(apiVersion)
}

func withResourceAPI(outputDir, apiVersion, kind string) string {
	return getOwnerAPIVersion(apiVersion, kind) + ` "` + getAppTypeImport(outputDir, apiVersion) + `"`
}

func getInstantiationString(resourceType string) string {
	base := resourceType[1:]
	return "&" + base + "{}"
}

func getImportMap(outputDir, kind, apiVersion string) map[string]string {

	controllerKindImports[getAppTypeImport(outputDir, apiVersion)] = getOwnerAPIVersion(apiVersion, kind)
	return controllerKindImports
}

func getAppTypeImport(outputDir, apiVersion string) string {
	var comps []string
	components, err := getAPIVersionComponents(apiVersion)
	if err != nil {
		panic(err)
	}

	sp := strings.Split(outputDir, "src/")
	comps = append([]string{sp[len(sp)-1], "pkg", "apis"}, components.Subdomain, components.Version)

	return filepath.Join(comps...)

}

func getResourceCRImport(outputDir, resourceType string) string {
	// get the beginning part
	sp := strings.Split(outputDir, "src/")
	// add on pkg, resources
	importString := sp[len(sp)-1] + "/pkg" + "/resources"
	// add on the resource type
	importString = importString + "/" + getNamePlural(resourceType)
	// return
	return importString
}

func getAPIVersionComponents(input string) (*APIComponents, error) {

	// matches the input string and returns the groups
	pattern := regexp.MustCompile(`(.*)\.(.*)\..*\/(.*)`)
	matches := pattern.FindStringSubmatch(input)
	if l := len(matches); l != 3+1 {
		return nil, fmt.Errorf("expected four matches, received %d instead", l)
	}

	var result = &APIComponents{
		matches[1],
		matches[2],
		matches[3],
	}

	return result, nil
}

// APIComponents used to seperate api components
type APIComponents struct {
	Subdomain string
	Domain    string
	Version   string
}

func getKindSpecName(kind string) string {
	return kind + "Spec"
}

func getKindSpecTypeName(kind string) string {
	return kind + "Spec" + "Type"
}

func getLowerPlural(kind string) string {
	return inflect.Pluralize(inflect.CamelizeDownFirst(kind))
}

func getSpecifiedOrDefault(attribute string, values map[string]string) string {
	// check if attribute is in the map
	val, ok := values[attribute]
	if !ok {
		return getDefault(attribute)
	}
	return val
}

func getDefault(attribute string) string {
	return "placeholder"
}
