package templating

import (
	"testing"
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

	outputDir = "/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/nginx-operator"
	apiVersion = "web.example.com/v1alpha1"
	expected = "github.com/redhat-nfvpe/helm2go-operator-sdk/nginx-operator/web/example/v1alpha1"
	result = getAppTypeImport(outputDir, apiVersion)
	if result != expected {
		t.Fatalf("expected %s got %s", expected, result)
	}
}
