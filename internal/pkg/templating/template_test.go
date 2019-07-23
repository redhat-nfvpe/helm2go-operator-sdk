package templating

import (
	"bytes"
	"testing"
	"text/template"
)

func TestLowerCamel(t *testing.T) {
	var input string
	var expected string
	var res string

	input = "nginx"
	expected = "nginx"
	if res = kindToLowerCamel(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}

	input = "TensorflowNotebook"
	expected = "tensorflowNotebook"
	if res = kindToLowerCamel(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}

	input = "Tensorflow Notebook"
	expected = "tensorflowNotebook"
	if res = kindToLowerCamel(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}
}

func TestLower(t *testing.T) {
	var input string
	var expected string
	var res string

	input = "nginx"
	expected = "nginx"
	if res = kindToLower(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}

	input = "TensorflowNotebook"
	expected = "tensorflownotebook"
	if res = kindToLower(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}

	input = "Tensorflow Notebook"
	expected = "tensorflow notebook"
	if res = kindToLower(input); res != expected {
		t.Fatalf("expected %s got %s", expected, res)
	}
}

func TestGetOwnerAPIVersion(t *testing.T) {
	var apiVersion string
	var kind string
	var res string

	apiVersion = "web.example.com/v1"
	kind = "nginx"

	if res = getOwnerAPIVersion(apiVersion, kind); res != "nginxv1" {
		t.Fatalf("expected 'nginxv1' got %s", res)
	}

	apiVersion = "apps.example.com/v1alpha1"
	kind = "TensorflowNotebook"

	if res = getOwnerAPIVersion(apiVersion, kind); res != "tensorflowNotebookv1alpha1" {
		t.Fatalf("expected 'tensorflowNotebookv1alpha1' got %s", res)
	}

}

func TestGetImport(t *testing.T) {
	var outputDir string
	var apiVersion string
	var result string
	var expected string

	outputDir = "/home/user/go/src/github.com/user/nginx-operator"
	apiVersion = "web.example.com/v1alpha1"
	expected = "github.com/user/nginx-operator/pkg/apis/web/v1alpha1"
	result = getAppTypeImport(outputDir, apiVersion)
	if result != expected {
		t.Fatalf("expected %s got %s", expected, result)
	}
}

func TestReconcileBlockRender(t *testing.T) {
	var result string
	var expected string

	reconcileConfig := ReconcileTemplateConfig{
		Kind:                    "Memcached",
		LowerKind:               "memcached",
		ResourceImportPackage:   "appsv1",
		ResourceType:            "Deployment",
		ResourceName:            "Deployment",
		LowerResourceName:       "deployment",
		PluralLowerResourceName: "deployments",
	}

	wr := bytes.Buffer{}

	temp, err := template.New("controllerFuncTemplate").Parse(reconcileConfig.GetTemplate())
	if err != nil {
		t.Fatalf("unexpected error while loading reconcile block template: %v", err)
	}
	// execute the template
	err = temp.Execute(&wr, reconcileConfig)
	if err != nil {
		t.Fatalf("unexpected error while rendering reconcile block template: %v", err)
	}
	result = wr.String()
	// set the expected value
	expected = `
		deployment := deployments.NewDeploymentForCR(instance)
		reqLogger.Info(deployment.String())
		// Set Memcached instance as the owner and controller
		if errDeployment := controllerutil.SetControllerReference(instance, deployment, r.scheme); errDeployment != nil {
			return reconcile.Result{}, errDeployment
		}
		// Check if this Deployment already exists
		foundDeployment := &appsv1.Deployment{}
		errDeployment := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: request.Namespace}, foundDeployment)
		if errDeployment != nil && errors.IsNotFound(errDeployment) {
			reqLogger.Info("Creating a new Deployment", "deployment.Namespace", request.Namespace, "deployment.Name", deployment.Name)
			deployment.ObjectMeta.SetNamespace(request.Namespace)
			errDeployment = r.client.Create(context.TODO(), deployment)
			if errDeployment != nil {
				return reconcile.Result{}, errDeployment
			}
			// Pod created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if errDeployment != nil {
			return reconcile.Result{}, errDeployment
		}
		reqLogger.Info("Skip reconcile: deployment already exists", "Deployment.Namespace",
		foundDeployment.Namespace, "svcacdeploymentcnt.Name", foundDeployment.Name)
	`

	if result != expected {
		t.Fatalf("uexpected error; expected reconcile result got: %s", result)
	}
}
