package render

import (
	"fmt"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/renderutil"
)

// InjectedToTemp ...
func InjectedToTemp(files map[string]string, outParentDir string) (string, error) {
	s, err := writeToTemp(files, outParentDir)
	return s, err
}

// InjectTemplateValues renders the template values from a valid Helm chart
// func InjectTemplateValues(repo bool, url bool, chart string) (map[string]string, error) {
func InjectTemplateValues(chartPath string) (map[string]string, error) {
	// var (
	// 	deleteTemp bool
	// 	tempDir string
	// )
	// // TODO parser to check if its a repo or url
	// 	// if it is, then download them and save them to a temp folder
	// if repo {
	// 	// fetch with the command and put in temp

	// 	deleteTemp = true
	// }
	// else if url {
	// 	// fetch with the command and put in temp
	// 	deleteTemp = true
	// }
	// load the raw chart; chartPath follows the specifications as listed in the helm documentation
	c, err := chartutil.Load(chartPath)
	if err != nil {
		fmt.Println("broken")
		return nil, err
	}
	// Inject the values from the raw chart
	rendered, err := renderutil.Render(c, c.GetValues(), renderutil.Options{})
	if err != nil {
		return nil, err
	}

	// if deleteTemp {
	// 	// delete the temp directory
	// }

	// return map; keys are existing filenames; values are file contents
	return rendered, nil
}
