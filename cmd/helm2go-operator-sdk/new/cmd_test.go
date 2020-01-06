package new

import (
	"fmt"
	"os"
	"path/filepath"

	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/pathconfig"
	"github.com/spf13/cobra"
)

var cwd, _ = os.Getwd()
var parent = filepath.Dir(filepath.Dir(filepath.Dir(cwd)))

func TestCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command")
}

var _ = Describe("Load Chart", func() {
	It("Gets Local Chart", func() {
		var local = "/test/bitcoind"
		var testLocal = parent + local
		chartClient := NewChartClient()
		chartClient.HelmChartRef = testLocal
		err := chartClient.LoadChart()
		Expect(err).ToNot(HaveOccurred())
		Expect(chartClient.Chart.GetMetadata().Name).To(Equal("bitcoind"))
	})
	It("Gets External Charts From Valid Repos", func() {
		var client HelmChartClient

		gopathDir := os.Getenv("GOPATH")
		pathConfig := pathconfig.NewConfig(filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk"))
		client.PathConfig = pathConfig

		client.HelmChartRef = "auto-deploy-app"
		client.HelmChartRepo = "https://charts.gitlab.io/"
		err := client.LoadChart()
		Expect(err).ToNot(HaveOccurred())

		// cleanup auto-deploy-app

		client.HelmChartRef = "not-a-chart"
		client.HelmChartRepo = ""
		err = client.LoadChart()
		Expect(err).To(HaveOccurred())

		client.HelmChartRepo = "invalidrepo.io"
		err = client.LoadChart()
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Flag Validation", func() {
	It("Verifies Command Flags", func() {
		var err error
		helmChartRef = ""
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		helmChartRef = "./test/tomcat"
		kind = "Tomcat"
		apiVersion = ""
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		apiVersion = "app.example/v1alpha1"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())
		apiVersion = "app.example.com/v1alpha1"
		err = verifyFlags()
		Expect(err).ToNot(HaveOccurred())
	})
	It("Verifies Kind Flag", func() {
		var err error
		apiVersion = "app.example.com/v1alpha1"
		helmChartRef = "./test/tomcat"

		kind = ""
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		kind = "tensorflow-notebook"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		kind = "tensorflowNotebook"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		kind = "Tensorflow-Notebook"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		kind = "TensorflowNotebook"
		err = verifyFlags()
		Expect(err).ToNot(HaveOccurred())
	})
	It("Verifies API Version", func() {
		var err error
		helmChartRef = "./test/tomcat"
		kind = "TensorflowNotebook"

		apiVersion = "app.example.com/v1alpha1"
		err = verifyFlags()
		Expect(err).ToNot(HaveOccurred())

		apiVersion = "web.k8s.io/v1"
		err = verifyFlags()
		Expect(err).ToNot(HaveOccurred())

		apiVersion = "k8s.io"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		apiVersion = "app,s.k8s.io/v1beta1"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

		apiVersion = "apps.example/v1beta1"
		err = verifyFlags()
		Expect(err).To(HaveOccurred())

	})
})

var _ = Describe("Command Argument Validation", func() {
	It("Fails On Invalid Arguments", func() {
		var err error
		err = createInvalidCommand().Execute()
		Expect(err).To(HaveOccurred())

		err = createInvalidOperatorName().Execute()
		Expect(err).To(HaveOccurred())

		err = createInvalidKindName().Execute()
		Expect(err).To(HaveOccurred())
	})
	It("Passes On Valid Arguments", func() {
		var err error
		err = createValidCommand().Execute()
		Expect(err).ToNot(HaveOccurred())
		cleanupValidOperator()
	})
})

//creates a command for use in the argument validation test
func createValidCommand() *cobra.Command {
	var cmd *cobra.Command
	cmd = GetNewCmd()

	gopathDir := os.Getenv("GOPATH")
	nginxDir := filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/nginx")

	// command has an operator name and correct flags
	cmd.SetArgs([]string{
		"test-operator",
		fmt.Sprintf("--helm-chart=%s", nginxDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Nginx"),
		fmt.Sprintf("--mock=%s", "true"),
	})
	return cmd
}

func cleanupValidOperator() {
	gopathDir := os.Getenv("GOPATH")
	operatorDir := filepath.Join(gopathDir, "src", "redhat-nfvpe", "helm2go-operator-sdk", "cmd", "helm2go-operator-sdk", "new", "test-operator")
	fmt.Println(operatorDir)
	os.RemoveAll(operatorDir)
}

// creates an invalid command for use in the argument validation test
func createInvalidCommand() *cobra.Command {
	var cmd *cobra.Command
	cmd = GetNewCmd()

	gopathDir := os.Getenv("GOPATH")
	nginxDir := filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/nginx")

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		fmt.Sprintf("--helm-chart=%s", nginxDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Ngnix"),
		fmt.Sprintf("--mock=%s", "true"),
	})
	return cmd
}

func createInvalidOperatorName() *cobra.Command {
	var cmd *cobra.Command
	cmd = GetNewCmd()

	gopathDir := os.Getenv("GOPATH")
	nginxDir := filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/nginx")

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		"",
		fmt.Sprintf("--helm-chart=%s", nginxDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Ngnix"),
		fmt.Sprintf("--mock=%s", "true"),
	})
	return cmd
}

func createInvalidKindName() *cobra.Command {
	var cmd *cobra.Command
	cmd = GetNewCmd()

	gopathDir := os.Getenv("GOPATH")
	nginxDir := filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/nginx")

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		"test-operator",
		fmt.Sprintf("--helm-chart=%s", nginxDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1 Server"),
		fmt.Sprintf("--kind=%s", "Nginx"),
	})
	return cmd
}
