package convert

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/repo"

	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/load"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/render"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/templating"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
)

func loadChart() error {
	// if repo is specified
	var chartPath string
	chartPath = helmChartRef
	if len(helmChartRepo) > 0 {

		var out io.Writer
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		d := downloader.ChartDownloader{
			Out:      out,
			Verify:   downloader.VerifyNever,
			Keyring:  "",
			HelmHome: helmpath.Home(filepath.Join(home, ".helm")),
			Getters:  getter.All(environment.EnvSettings{}),
			Username: username,
			Password: password,
		}

		chartURL, err := repo.FindChartInAuthRepoURL(helmChartRepo, username, password, chartPath, helmChartVersion, helmChartCertFile, helmChartKeyFile, helmChartCAFile, getter.All(environment.EnvSettings{}))
		if err != nil {
			return err
		}

		helmChartRef = chartURL
		cwd, _ := os.Getwd()

		downloaded, _, err := d.DownloadTo(helmChartRef, helmChartVersion, filepath.Join(cwd, outputDir))
		if err != nil {
			return err
		}

		ud := chartName
		if !filepath.IsAbs(ud) {
			ud = filepath.Join(outputDir, ud)
		}
		if fi, err := os.Stat(ud); err != nil {
			if err := os.MkdirAll(ud, 0755); err != nil {
				return fmt.Errorf("Failed to untar (mkdir): %s", err)
			}

		} else if !fi.IsDir() {
			return fmt.Errorf("Failed to untar: %s is not a directory", ud)
		}

		chartutil.ExpandFile(ud, downloaded)
		os.RemoveAll(downloaded)
		log.Printf("Downloaded Chart To: %v\n", ud)
		chartPath = filepath.Join(ud, chartPath)
	}

	c, _ = chartutil.Load(chartPath)
	chartName = c.Metadata.GetName()

	return nil
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