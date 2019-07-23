package render

import (
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
)

// InjectedToTemp ...
func InjectedToTemp(files map[string]string, outParentDir string) (string, error) {
	s, err := writeToTemp(files, outParentDir)
	return s, err
}

// InjectTemplateValues renders the template values from a valid Helm chart
func InjectTemplateValues(c *chart.Chart) (map[string]string, error) {
	// Inject the values from the raw chart
	rendered, err := renderutil.Render(c, c.GetValues(), renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name: "release",
		},
	})
	if err != nil {
		return nil, err
	}
	return rendered, nil
}
