package render

import (
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/renderutil"
)

// InjectTemplateValues renders the template values from a valid Helm chart
func InjectTemplateValues(chartPath string) (map[string]string, error) {

	// load the raw chart; chartPath follows the specifications as listed in the helm documentation
	c, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}
	// Inject the values from the raw chart
	rendered, err := renderutil.Render(c, c.GetValues(), renderutil.Options{})
	if err != nil {
		return nil, err
	}
	// return map; keys are existing filenames; values are file contents
	return rendered, nil
}
