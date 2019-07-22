package templating

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
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

func kindToLowerCamel(kind string) string {
	return strcase.ToLowerCamel(kind)
}

func getOwnerAPIVersion(apiVersion, kind string) string {
	return kindToLowerCamel(kind) + filepath.Base(apiVersion)
}

func getImportMap(outputDir, kind, apiVersion string) map[string]string {

	controllerKindImports[getAppTypeImport(outputDir, apiVersion)] = getOwnerAPIVersion(apiVersion, kind)
	return controllerKindImports
}

func getAppTypeImport(outputDir, apiVersion string) string {
	comps, err := getAPIVersionComponents(apiVersion)
	if err != nil {
		panic(err)
	}
	sp := strings.Split(outputDir, "src/")
	comps = append([]string{sp[len(sp)-1], "pkg", "apis"}, comps...)

	return filepath.Join(comps...)

}

func getAPIVersionComponents(input string) ([]string, error) {

	var group string
	var version string

	// matches the input string and returns the groups
	pattern := regexp.MustCompile(`(.*)\.(.*)\..*\/(.*)`)
	matches := pattern.FindStringSubmatch(input)
	if l := len(matches); l != 3+1 {
		return []string{}, fmt.Errorf("expected four matches, received %d instead", l)
	}
	group = matches[1]
	version = matches[3]

	var result = []string{
		group,
		version,
	}

	return result, nil
}

func getControllerTemplate() string {
	return `
		package {{ .LowerKind }}
		import (
			"context"
			{{range $p, $i := .ImportMap -}}
			{{$i}} "{{$p}}"
			{{end}}
		)
		var log = logf.Log.WithName("controller_{{ .LowerKind }}")
		/**
		* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
		* business logic.  Delete these comments after modifying this file.*
		*/
		// Add creates a new {{ .Kind }} Controller and adds it to the Manager. The Manager will set fields on the Controller
		// and Start it when the Manager is Started.
		func Add(mgr manager.Manager) error {
			return add(mgr, newReconciler(mgr))
		}
		// newReconciler returns a new reconcile.Reconciler
		func newReconciler(mgr manager.Manager) reconcile.Reconciler {
			return &Reconcile{{ .Kind }}{client: mgr.GetClient(), scheme: mgr.GetScheme()}
		}
		// add adds a new Controller to mgr with r as the reconcile.Reconciler
		func add(mgr manager.Manager, r reconcile.Reconciler) error {
			// Create a new controller
			c, err := controller.New("{{ .LowerKind }}-controller", mgr, controller.Options{Reconciler: r})
			if err != nil {
				return err
			}
			// Watch for changes to primary resource {{ .Kind }}
			err = c.Watch(&source.Kind{Type: &{{ .OwnerAPIVersion }}.{{ .Kind }}{}}, &handler.EnqueueRequestForObject{})
			if err != nil {
				return err
			}
			// TODO(user): Modify this to be the types you create that are owned by the primary resource
			// Watch for changes to secondary resource Pods and requeue the owner {{ .Kind }}
			err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &{{ .OwnerAPIVersion }}.{{ .Kind }}{},
			})
			if err != nil {
				return err
			}

			// GENERATED BY CONVERSION KIT

			{{range $f :=  .ResourceWatches -}}
				{{$f}}
				
			{{end}}

			return nil
		}
		// blank assignment to verify that Reconcile{{ .Kind }} implements reconcile.Reconciler
		var _ reconcile.Reconciler = &Reconcile{{ .Kind }}{}
		// Reconcile{{ .Kind }} reconciles a {{ .Kind }} object
		type Reconcile{{ .Kind }} struct {
			// This client, initialized using mgr.Client() above, is a split client
			// that reads objects from the cache and writes to the apiserver
			client client.Client
			scheme *runtime.Scheme
		}
		// Reconcile reads that state of the cluster for a {{ .Kind }} object and makes changes based on the state read
		// and what is in the {{ .Kind }}.Spec
		// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
		// a Pod as an example
		// Note:
		// The Controller will requeue the Request to be processed again if the returned error is non-nil or
		// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
		func (r *Reconcile{{ .Kind }}) Reconcile(request reconcile.Request) (reconcile.Result, error) {
			reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
			reqLogger.Info("Reconciling {{ .Kind }}")
			// Fetch the {{ .Kind }} instance
			instance := &{{ .OwnerAPIVersion }}.{{ .Kind }}{}
			err := r.client.Get(context.TODO(), request.NamespacedName, instance)
			if err != nil {
				if errors.IsNotFound(err) {
					// Request object not found, could have been deleted after reconcile request.
					// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
					// Return and don't requeue
					return reconcile.Result{}, nil
				}
				// Error reading the object - requeue the request.
				return reconcile.Result{}, err
			}
			// Define a new Pod object
			pod := newPodForCR(instance)
			// Set {{ .Kind }} instance as the owner and controller
			if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
				return reconcile.Result{}, err
			}
			// Check if this Pod already exists
			found := &corev1.Pod{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
			if err != nil && errors.IsNotFound(err) {
				reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
				err = r.client.Create(context.TODO(), pod)
				if err != nil {
					return reconcile.Result{}, err
				}
				// Pod created successfully - don't requeue
				return reconcile.Result{}, nil
			} else if err != nil {
				return reconcile.Result{}, err
			}
			// Pod already exists - don't requeue
			reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
			return reconcile.Result{}, nil
		}
		// newPodForCR returns a busybox pod with the same name/namespace as the cr
		func newPodForCR(cr *{{ .OwnerAPIVersion }}.{{ .Kind }}) *corev1.Pod {
			labels := map[string]string{
				"app": cr.Name,
			}
			return &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name + "-pod",
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "busybox",
							Image:   "busybox",
							Command: []string{"sleep", "3600"},
						},
					},
				},
			}
		}
`
}

func getControllerFunctionTemplate() string {
	return `
	err = c.Watch(&source.Kind{Type: &{{.ResourceImportPackage}}.{{.ResourceType}}{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &{{.OwnerAPIVersion}}.{{.Kind}}{},
	})
	if err != nil {
		return err
	}
	`
}

func getResourceTemplate() string {
	return `
	package {{ .PackageName }}

	import (
		{{range $index, $statement := .ImportStatements}}
			{{ $statement }}
		{{ end }}
	)

	// New{{ .ResourceName }}ForCR ...
	func New{{ .ResourceName }}ForCR(r *{{.APIVersion}}.{{ .Kind }}) {{ .ResourceType }}{
		var e {{ .ResourceType }}
		elemYaml := ` + "`" + `{{ .ResourceJSON }}` + "`" + `
		// Unmarshal Specified JSON to Kubernetes Resource
		err := json.Unmarshal([]byte(elemYaml), e)
		if err != nil {
			panic(err)
		}
		return e
	}
	`
}
