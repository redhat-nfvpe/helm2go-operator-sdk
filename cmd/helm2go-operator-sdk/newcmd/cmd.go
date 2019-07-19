package newcmd

import (
	"path/filepath"

	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/pathconfig"
	"github.com/spf13/cobra"
	"github.com/tav/golly/log"
)

// NewCmd ...
func NewCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new <New Name>",
		Short: "Creates a Go Operator from an Existing Helm Chart",
		Long:  "Utilizes the Helm Rendering Engine and Operator-SDK to consume an existing Helm Chart and produces a Go Operator",
		RunE:  convertFunc,
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
)

var (
	c          *chart.Chart
	chartName  string
	pathConfig *pathconfig.PathConfig
)

func convertFunc(cmd *cobra.Command, args []string) error {

	setBasePathConfig()

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

	log.Infof("🤠 Converting Existing Helm Chart %s to Go Operator %s!", helmChartRef, operatorName)

	// load the spcecified helm chart
	err := loadChart()
	if err != nil {
		log.Error("error loading chart: ", err)
		return err
	}

	rcache, err := doHelmGoConversion()
	if err != nil {
		log.Error("error performing chart conversion: ", err)
		return err
	}

	outputDir = filepath.Join(pathConfig.GetBasePath(), operatorName)

	//create the operator-sdk scaffold
	_, err = doGoScaffold()
	if err != nil {
		log.Error("error generating scaffolding: ", err)
		return err
	}

	err = scaffoldOverwrite(outputDir, kind, apiVersion, rcache)
	if err != nil {
		log.Error("error overwritting scaffold: ", err)
		return err
	}

	log.Infof("🤠 Go Operator Can Be Found At: %s", outputDir)
	return nil
}
