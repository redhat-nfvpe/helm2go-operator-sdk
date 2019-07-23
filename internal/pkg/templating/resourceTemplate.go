package templating

import (
	"bytes"
	"reflect"
	"text/template"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// ResourceTemplateConfig is used as payload for executing the templating functionality
type ResourceTemplateConfig struct {
	PackageName               string
	Kind                      string
	APIVersion                string
	OwnerAPIVersion           string
	ResourceName              string
	ResourceType              string
	ResourceTypeInstantiation string
	ResourceJSON              string
	ImportStatements          []string
	Resource                  interface{}
}

// NewResourceTemplateConfig returns a new resource templating configuration object
func NewResourceTemplateConfig(outputDir, apiVersion, kind string, r *resourcecache.Resource, bytes []byte) *ResourceTemplateConfig {
	return &ResourceTemplateConfig{
		ImportStatements:          append(getTemplateImports(&r.GetResourceFunctions()[0].Data), withResourceAPI(outputDir, apiVersion, kind)),
		PackageName:               string(r.PackageName),
		Kind:                      kind,
		ResourceName:              getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
		ResourceType:              reflect.TypeOf((r.GetResourceFunctions()[0].Data)).String(),
		ResourceTypeInstantiation: getInstantiationString(reflect.TypeOf((r.GetResourceFunctions()[0].Data)).String()),
		Resource:                  getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
		ResourceJSON:              string(bytes),
		APIVersion:                apiVersion,
		OwnerAPIVersion:           getOwnerAPIVersion(apiVersion, kind),
	}
}

// Execute renders the template and returns the templated string
func (r *ResourceTemplateConfig) Execute() (string, error) {
	temp, err := template.New("resourceFuncTemplate").Parse(r.GetTemplate())
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer
	err = temp.Execute(&wr, r)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

// GetTemplate returns the necessary template
func (r *ResourceTemplateConfig) GetTemplate() string {
	return `
		package {{ .PackageName }}

		import (
			{{range $index, $statement := .ImportStatements}}
				{{ $statement }}
			{{ end }}
		)

		// New{{ .ResourceName }}ForCR ...
		func New{{ .ResourceName }}ForCR(r *{{.OwnerAPIVersion}}.{{ .Kind }}) {{ .ResourceType }}{
			e :=  {{ .ResourceTypeInstantiation }}
			elemJSON := ` + "`" + `{{ .ResourceJSON }}` + "`" + `
			// Unmarshal Specified JSON to Kubernetes Resource
			err := json.Unmarshal([]byte(elemJSON), e)
			if err != nil {
				panic(err)
			}
			return e
		}
	`
}
