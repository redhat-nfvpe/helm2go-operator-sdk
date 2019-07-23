package templating

import (
	"bytes"
	"text/template"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// ControllerWatchFuncConfig contains the necessary fields to render the controller watch function specific to a resource
type ControllerWatchFuncConfig struct {
	APIVersion            string
	ResourceImportPackage string
	ResourceType          string
	Kind                  string
	LowerKind             string
	OwnerAPIVersion       string
}

// NewControllerWatchFuncConfig returns a new controller watch function configuration object
func NewControllerWatchFuncConfig(apiVersion, ownerAPIVersion, kind, lowerKind string, r *resourcecache.Resource) *ControllerWatchFuncConfig {
	resourceType := getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)
	return &ControllerWatchFuncConfig{
		APIVersion:            apiVersion,
		OwnerAPIVersion:       ownerAPIVersion,
		Kind:                  kind,
		LowerKind:             lowerKind,
		ResourceImportPackage: getResourceImportPackage(resourceType),
		ResourceType:          resourceType,
	}
}

// Execute renders the template and returns the templated string
func (c *ControllerWatchFuncConfig) Execute() (string, error) {
	temp, err := template.New("resourceFuncTemplate").Parse(c.GetTemplate())
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

// GetTemplate returns the necessary template
func (c *ControllerWatchFuncConfig) GetTemplate() string {
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
