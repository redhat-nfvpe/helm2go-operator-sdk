package new

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tav/golly/log"
)

// GetNewCmd ...
func GetNewCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new <New Name>",
		Short: "Builds a Go Operator from an existing Helm Chart",
		Long:  "Utilizes the Helm Rendering Engine and Operator-SDK to consume an existing Helm Chart to produce a Go Operator",
		RunE:  newFunc,
	}

	newCmd.Flags().StringVar(&helmChartRef, "helm-chart", "", "Initialize helm operator with existing helm chart (<URL>, <repo>/<name>, or local path)")
	newCmd.Flags().StringVar(&helmChartVersion, "helm-chart-version", "", "Specific version of the helm chart (default is latest version)")
	newCmd.Flags().StringVar(&helmChartRepo, "helm-chart-repo", "", "Chart repository URL for the requested helm chart")
	newCmd.Flags().StringVar(&username, "username", "", "Username for chart repo")
	newCmd.Flags().StringVar(&password, "password", "", "Password for chart repo")
	newCmd.Flags().StringVar(&helmChartCertFile, "helm-chart-cert-file", "", "Cert File For External Repo")
	newCmd.Flags().StringVar(&helmChartKeyFile, "helm-chart-key-file", "", "Key File For External Repo")
	newCmd.Flags().StringVar(&helmChartCAFile, "helm-chart-ca-file", "", "CA File For External Repo")
	newCmd.Flags().StringVar(&apiVersion, "api-version", "", "Kubernetes apiVersion and has a format of $GROUP_NAME/$VERSION (e.g app.example.com/v1alpha1)")
	newCmd.Flags().StringVar(&kind, "kind", "", "Kubernetes CustomResourceDefintion kind. (e.g AppService)")
	newCmd.Flags().BoolVar(&clusterScoped, "cluster-scoped", false, "Operator cluster scoped or not")

	// debug flags
	newCmd.Flags().BoolVar(&mock, "mock", false, "Used for testing")
	// newCmd.Flags().MarkHidden("mock")

	return newCmd
}

var (
	helmChartRef      string
	helmChartVersion  string
	helmChartRepo     string
	username          string
	password          string
	helmChartCertFile string
	helmChartKeyFile  string
	helmChartCAFile   string
	apiVersion        string
	kind              string
	clusterScoped     bool
	outputDir         string
	operatorName      string
	mock              bool
)

func newFunc(cmd *cobra.Command, args []string) error {

	chartClient := NewChartClient()
	chartClient.SetValues(helmChartRef, helmChartVersion, helmChartRepo, username, password, helmChartCAFile, helmChartCertFile, helmChartKeyFile)

	if err := parse(args); err != nil {
		log.Error("error parsing arguments: ", err)
		return err
	}
	if err := verifyFlags(); err != nil {
		log.Error("error verifying flags: ", err)
		return err
	}

	if err := verifyOperatorSDKVersion(); err != nil {
		log.Error("error verifying operator-sdk version: ", err)
		return err
	}

	// for testing
	if mock {
		return nil
	}

	log.Infof("🤠 Creating Go Operator %s from Helm Chart %s!", operatorName, chartClient.HelmChartRef)

	// load the spcecified helm chart
	err := chartClient.LoadChart()

	if err != nil {
		log.Error("error loading chart: ", err)
		return err
	}

	// convert helm resources to Go resource cache
	rcache, err := chartClient.DoHelmGoConversion()
	if err != nil {
		log.Error("error performing chart conversion: ", err)
		return err
	}

	// output directory is the path/to/command/operator-name
	outputDir = filepath.Join(chartClient.PathConfig.GetBasePath(), operatorName)

	//create the operator-sdk scaffold
	err = doGoScaffold()
	if err != nil {
		log.Error("error generating scaffolding: ", err)
		return err
	}

	valsPath := filepath.Join(chartClient.ChartPath, "values.yaml")
	err = scaffoldOverwrite(outputDir, kind, apiVersion, valsPath, rcache)
	if err != nil {
		log.Error("error overwritting scaffold: ", err)
		return err
	}

	log.Infof("🤠 Go Operator Can Be Found At: %s", outputDir)
	return nil
}
