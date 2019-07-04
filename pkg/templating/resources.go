package templating

import (
	"path/filepath"
	"reflect"
	"strings"

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

func kindToLowerCamel(kind string) string {
	return strcase.ToLowerCamel(kind)
}

func getOwnerAPIVersion(apiVersion, kind string) string {
	return filepath.Base(apiVersion) + kindToLowerCamel(kind)
}

func getImportMap(outputDir, kind, apiVersion string) map[string]string {
	controllerKindImports[getAppTypeImport(outputDir, apiVersion)] = getAppTypeImportAbbreviation(kind, apiVersion)
	return controllerKindImports
}

func getAppTypeImport(outputDir, apiVersion string) string {
	// append the correct path
	// everything after source is in the correct path
	sp := strings.Split(outputDir, "src/")
	base := sp[len(sp)-1]
	result := filepath.Join(base, "pkg", "apis", "apps", apiVersion)
	return result
}

func getAppTypeImportAbbreviation(kind, apiVersion string) string {
	result := apiVersion + kind
	return result
}
