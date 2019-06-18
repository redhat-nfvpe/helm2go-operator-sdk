package convert

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/load"
	"github.com/redhat-nfvpe/helm2go-operator-sdk/pkg/templating"
	"github.com/spf13/cobra"
	"github.com/tav/golly/log"
)

// NewConvertCmd ...
func NewConvertCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "convert",
		Short: "Converts an existing Helm Chart Operator in to a Go Operator",
		Long:  "Utilizes the Helm Rendering Engine and Operator-SDK to consumer an existing Helm Chart operator and then produces a Go Operator",
		RunE:  convertFunc,
	}

	newCmd.Flags().StringVar(&helmChartRef, "helm-chart", "", "Initialize helm operator with existing helm chart (<URL>, <repo>/<name>, or local path)")
	newCmd.Flags().StringVar(&helmChartVersion, "helm-chart-version", "", "Specific version of the helm chart (default is latest version)")
	newCmd.Flags().StringVar(&helmChartRepo, "helm-chart-repo", "", "Chart repository URL for the requested helm chart")
	newCmd.Flags().StringVar(&apiVersion, "api-version", "", "Kubernetes apiVersion and has a format of $GROUP_NAME/$VERSION (e.g app.example.com/v1alpha1)")
	newCmd.Flags().StringVar(&kind, "kind", "", "Kubernetes CustomResourceDefintion kind. (e.g AppService)")
	newCmd.Flags().BoolVar(&clusterScoped, "cluster-scoped", false, "Operator cluster scoped or not")

	return newCmd
}

var (
	helmChartRef     string
	helmChartVersion string
	helmChartRepo    string
	apiVersion       string
	kind             string
	clusterScoped    bool
)

func convertFunc(cmd *cobra.Command, args []string) error {
	if err := parse(cmd, args); err != nil {
		return err
	}
	if err := verifyFlags(); err != nil {
		return err
	}

	log.Infof("🤠 Converting Existing Helm Chart %s to Go Operator!", helmChartRef)

	dir, err := doHelmGoConversion()
	if err != nil {
		return err
	}
	log.Infof("🤠 Go Operator Can Be Found At: %s", dir)
	return nil
}

func verifyFlags() error {
	// PLACEHOLDER
	return nil
}

func parse(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Please Use Specified Flags")
	}
	return nil
}

func doHelmGoConversion() (string, error) {

	fmt.Println(helmChartRef)

	// render the helm charts
	f, err := render.InjectTemplateValues(helmChartRef)
	if err != nil {
		return "", err
	}
	// write the rendered charts to output directory
	d, _ := os.Getwd()
	to, err := render.InjectedToTemp(f, d)
	if err != nil {
		return "", err
	}

	// convert the helm templates to go structures

	//TODO more robust filepath matching for the actual template directory
	to = filepath.Join(to, "drone", "templates")
	// cwd, _ := os.Getwd()
	// to := filepath.Join(cwd, "test", "resources")

	rcache, err := load.YAMLUnmarshalResources(to, resourcecache.NewResourceCache())
	if err != nil {
		return "", err
	}
	fmt.Printf("%v\n", rcache)

	templates := templating.CacheTemplating(rcache, "*v1alpha1.Collectd")
	fmt.Printf("%v\n", templates)
	// build the operator scaffold
	// pass the go structures cache
	return "", nil
}
