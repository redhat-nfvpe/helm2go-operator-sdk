package templating

import (
	"bytes"
	"text/template"
)

// CRTemplateConfig is used to template the appropriate CR Scaffold
type CRTemplateConfig struct {
	APIVersion string
	Kind       string
	LowerKind  string
	Spec       CRSpec
}

// CRSpec contains the specific values needed for the CRTemplate
type CRSpec struct {
	placeholder string
}

// NewCRTemplateConfig ...
// TODO : need to implement the function to create an actual spec struct
func NewCRTemplateConfig(apiVersion, kind, lowerKind string, spec CRSpec) *CRTemplateConfig {
	return &CRTemplateConfig{
		APIVersion: apiVersion,
		Kind:       kind,
		LowerKind:  lowerKind,
		Spec:       spec,
	}
}

// Execute renders the template and returns the templated string
func (c *CRTemplateConfig) Execute() (string, error) {
	temp, err := template.New("CRTemplate").Parse(c.GetTemplate())
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer
	err = temp.Execute(&wr, c)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

// GetTemplate returns the necessary template for the CR
func (c *CRTemplateConfig) GetTemplate() string {

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
