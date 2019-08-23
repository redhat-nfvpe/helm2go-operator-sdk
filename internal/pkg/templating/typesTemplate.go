package templating

// TypesTemplateConfig is used to render TypesTemplate
type TypesTemplateConfig struct {
	APIVersion string
	Kind       string
	LowerKind  string
	Spec       TypesSpec
}

// TypesSpec is the required spec for TypesTemplate
type TypesSpec struct {
	placeholder string
}

// GetTemplate returns the necessary template for the CR
func (t *TypesTemplateConfig) GetTemplate() string {
	return `apiVersion: {{ .Resource.APIVersion }}
	kind: {{ .Resource.Kind }}
	metadata:
	  name: example-{{ .Resource.LowerKind }}
	spec:
	{{- with .Spec }}
	{{ . | indent 2 }}
	{{- else }}
	  # Add fields here
	  size: 3
	{{- end }}
	`
}
