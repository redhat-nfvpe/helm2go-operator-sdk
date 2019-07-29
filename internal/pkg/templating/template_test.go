package templating

import (
	"bytes"
	"testing"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTemplating(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Templating")
}

var _ = Describe("LowerCamel", func() {
	var input string
	var expected string
	var res string

	BeforeEach(func() {
		input = ""
		expected = ""
		res = ""
	})

	It("Does Not Change Already Lower Camel Strings", func() {
		input = "nginx"
		expected = "nginx"
		res = kindToLowerCamel(input)
		Expect(res).To(Equal(expected))
	})
	It("Changes Upper Camel Strings", func() {
		input = "TensorflowNotebook"
		expected = "tensorflowNotebook"
		res = kindToLowerCamel(input)
		Expect(res).To(Equal(expected))

	})
	It("Removes Strings", func() {
		input = "Tensorflow Notebook"
		expected = "tensorflowNotebook"
		res = kindToLowerCamel(input)
		Expect(res).To(Equal(expected))
	})
})

var _ = Describe("Lower", func() {
	var input string
	var expected string
	var res string

	BeforeEach(func() {
		input = ""
		expected = ""
		res = ""
	})

	It("Does Not Change Lower Strings", func() {
		input = "nginx"
		expected = "nginx"
		res = kindToLower(input)
		Expect(res).To(Equal(expected))
	})
	It("Changes Upper Camel To Lower", func() {
		input = "TensorflowNotebook"
		expected = "tensorflownotebook"
		res = kindToLower(input)
		Expect(res).To(Equal(expected))
	})
	It("Removes Strings", func() {
		input = "Tensorflow Notebook"
		expected = "tensorflow notebook"
		res = kindToLower(input)
		Expect(res).To(Equal(expected))
	})
})

var _ = Describe("OwnerAPIVersion", func() {
	var apiVersion string
	var kind string
	var res string

	It("Gets The Correct API Version", func() {
		apiVersion = "web.example.com/v1"
		kind = "nginx"
		res = getOwnerAPIVersion(apiVersion, kind)
		Expect(res).To(Equal("nginxv1"))

		apiVersion = "apps.example.com/v1alpha1"
		kind = "TensorflowNotebook"
		res = getOwnerAPIVersion(apiVersion, kind)
		Expect(res).To(Equal("tensorflowNotebookv1alpha1"))
	})
})

var _ = Describe("GetImport", func() {
	var outputDir string
	var apiVersion string
	var result string
	var expected string

	It("Gets The Right Import Statement", func() {
		outputDir = "/home/user/go/src/github.com/user/nginx-operator"
		apiVersion = "web.example.com/v1alpha1"
		expected = "github.com/user/nginx-operator/pkg/apis/web/v1alpha1"
		result = getAppTypeImport(outputDir, apiVersion)
		Expect(result).To(Equal(expected))
	})
})

var _ = Describe("Reconcile Render", func() {
	var result string
	var expected string

	It("Outputs The Correct Values", func() {
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
		Expect(err).ToNot(HaveOccurred())
		// execute the template
		err = temp.Execute(&wr, reconcileConfig)
		Expect(err).ToNot(HaveOccurred())

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

		Expect(result).To(Equal(expected))
	})
})
