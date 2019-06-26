package convert

import (
	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/spf13/cobra"
	"github.com/tav/golly/log"
)

// NewConvertCmd ...
func NewConvertCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "convert <New Name>",
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
	outputDir        string
)

var (
	c         *chart.Chart
	chartName string
)

func convertFunc(cmd *cobra.Command, args []string) error {
	if err := parse(args); err != nil {
		log.Error(err)
		return err
	}
	if err := verifyFlags(); err != nil {
		log.Error(err)
		return err
	}

	log.Infof("🤠 Converting Existing Helm Chart %s to Go Operator %s!", helmChartRef, outputDir)

	// load the spcecified helm chart
	loadChart()

	// create the operator-sdk scaffold
	_, err := doGoScaffold()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = doHelmGoConversion()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("🤠 Go Operator Can Be Found At: %s", outputDir)
	return nil
}
