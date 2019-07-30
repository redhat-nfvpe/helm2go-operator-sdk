package templating

import (
	"bytes"
	"text/template"
)

// KindTypesTemplateConfig is used to overwrite the generated kind_types.go file
type KindTypesTemplateConfig struct {
	Kind                            string
	Version                         string
	LowerPlural                     string
	KindSpec                        string
	SpecConfigurationName           string
	SpecConfigurationTypeName       string
	LowerCamelSpecConfigurationName string
}

// NewKindTypesTemplateConfig returns a new KindTypesConfig
func NewKindTypesTemplateConfig(kind, apiVersion, kindSpec string) *KindTypesTemplateConfig {

	lowerPlural := getLowerPlural(kind)
	components, _ := getAPIVersionComponents(apiVersion)
	version := components.Version
	kindSpecName := getKindSpecName(kind)
	kindSpecTypeName := getKindSpecTypeName(kind)
	kindSpecLowerCamel := kindToLowerCamel(kind)

	// create config and return address
	return &KindTypesTemplateConfig{
		Kind:                            kind,
		Version:                         version,
		LowerPlural:                     lowerPlural,
		KindSpec:                        kindSpec,
		SpecConfigurationName:           kindSpecName,
		SpecConfigurationTypeName:       kindSpecTypeName,
		LowerCamelSpecConfigurationName: kindSpecLowerCamel,
	}
}

// Execute renders the template and returns the templated string
func (k *KindTypesTemplateConfig) Execute() (string, error) {
	temp, err := template.New("CRTemplate").Parse(k.GetTemplate())
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer
	err = temp.Execute(&wr, k)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

// GetTemplate returns the necessary template for the kind_types.go
func (k *KindTypesTemplateConfig) GetTemplate() string {
	return `package {{ .Version }}
	import (
		metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	)
	// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
	// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

	{{ .KindSpec }}
	

	// {{.Kind}}Spec defines the desired state of {{.Kind}}
	// +k8s:openapi-gen=true
	type {{.Kind}}Spec struct {
		// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
		// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
		// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

		{{ .SpecConfigurationName }} {{ .SpecConfigurationTypeName }}` + "`" + `json:"{{ .LowerCamelSpecConfigurationName }},omitempty"` + "`" + `
		
	}
	// {{.Kind}}Status defines the observed state of {{.Kind}}
	// +k8s:openapi-gen=true
	type {{.Kind}}Status struct {
		// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
		// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
		// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	}
	// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
	// {{.Kind}} is the Schema for the {{ .LowerPlural }} API
	// +k8s:openapi-gen=true
	// +kubebuilder:subresource:status
	type {{.Kind}} struct {
		metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
		metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
		Spec   {{.Kind}}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
		Status {{.Kind}}Status ` + "`" + `json:"status,omitempty"` + "`" + `
	}
	// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
	// {{.Kind}}List contains a list of {{.Kind}}
	type {{.Kind}}List struct {
		metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
		metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
		Items           []{{ .Kind }} ` + "`" + `json:"items"` + "`" + `
	}
	func init() {
		SchemeBuilder.Register(&{{.Kind}}{}, &{{.Kind}}List{})
	}
	`
}
