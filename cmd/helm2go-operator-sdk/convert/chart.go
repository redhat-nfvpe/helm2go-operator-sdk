package convert

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/load"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/render"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/templating"

	"k8s.io/helm/pkg/chartutil"
)

func loadChart() {
	// TODO add the external chart functionality
	c, _ = chartutil.Load(helmChartRef)
	if c == nil {
		log.Println("Chart Is Empty! Exiting")
	}
	chartName = c.Metadata.GetName()
}

func doHelmGoConversion() (string, error) {
	// render the helm charts
	f, err := render.InjectTemplateValues(c)
	if err != nil {
		return "", err
	}
	// write the rendered charts to output directory
	d, _ := os.Getwd()
	temp, err := render.InjectedToTemp(f, d)
	if err != nil {
		return "", err
	}

	// convert the helm templates to go structures
	to := filepath.Join(temp, chartName, "templates")

	rcache, err := load.YAMLUnmarshalResources(to, resourcecache.NewResourceCache())
	if err != nil {
		return "", err
	}

	// clean up temp folder
	os.RemoveAll(temp)

	ok := templating.OverwriteController(outputDir, kind, apiVersion, rcache)
	if !ok {
		fmt.Println(ok)
	}
	// create templates for writing to file
	templates := templating.CacheTemplating(rcache, kind, apiVersion)
	// templates to files; outputDir is the parent directory where the operator scaffolding lives
	resDir := filepath.Join(outputDir, "pkg", "resources")
	// create the necessary package resource specific folders
	ok = templating.ResourceFileStructure(rcache, resDir)
	ok = templating.TemplatesToFiles(templates, resDir)
	if !ok {
		return "", fmt.Errorf("Writing to File Error")
	}
	return "", nil
}
