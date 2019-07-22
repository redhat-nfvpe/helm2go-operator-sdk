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

func TestResourceTemplate(t *testing.T) {

	type ResourceTemplate struct {
		PackageName      string
		ImportStatements []string
		ResourceName     string
		APIVersion       string
		Kind             string
		ResourceType     string
		ResourceJSON     string
	}
	r := ResourceTemplate{}
	r.PackageName = "MyPackage"
	r.ImportStatements = []string{"A", "b", "c", "d", "e"}
	r.ResourceName = "Resource"
	r.APIVersion = "alpha"
	r.Kind = "Test"
	r.ResourceType = "RType"

	r.ResourceJSON = `{
		"glossary": {
			"title": "example glossary",
			"GlossDiv": {
				"title": "S",
				"GlossList": {
					"GlossEntry": {
						"ID": "SGML",
						"SortAs": "SGML",
						"GlossTerm": "Standard Generalized Markup Language",
						"Acronym": "SGML",
						"Abbrev": "ISO 8879:1986",
						"GlossDef": {
							"para": "A meta-markup language, used to create markup languages such as DocBook.",
							"GlossSeeAlso": ["GML", "XML"]
						},
						"GlossSee": "markup"
					}
				}
			}
		}
	}`

	x := `
		package {{ .PackageName }}
	
		import (
			{{range $index, $statement := .ImportStatements}}
				{{ $statement }}
			{{ end }}
		)
	
		// New{{ .ResourceName }}ForCR ...
		func New{{ .ResourceName }}ForCR(r *{{.APIVersion}}.{{ .Kind }}) {{ .ResourceType }}{
			var e {{ .ResourceType }}
			elemYaml :=` + "`" + "{{ .ResourceJSON }}" + "`" + `
			// Unmarshal Specified JSON to Kubernetes Resource
			err := json.Unmarshal([]byte(elemYaml), e)
			if err != nil {
				panic(err)
			}
			return e
		}
		`
	tmpl, err := template.New("test").Parse(x)
	if err != nil {
		panic(err)
	}
	var s bytes.Buffer
	err = tmpl.Execute(&s, r)
	if err != nil {
		panic(err)
	}
}
