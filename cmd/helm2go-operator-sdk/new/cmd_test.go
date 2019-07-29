package new

import (
	"fmt"
	"io/ioutil"
	"log"
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

func silencelog() {
	log.SetOutput(ioutil.Discard)
}
func resetlog() {
	log.SetOutput(os.Stdout)
}

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

func TestInvalidCommandArgumentValidation(t *testing.T) {
	var err error
	silencelog()
	if err = createInvalidCommand().Execute(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Expected Error! Command Has No Operator Name Argument.")
		silencelog()
	}
	if err = createInvalidOperatorName().Execute(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Expected Error! Command Has Invalid(empty) Operator Name Argument.")
		silencelog()
	}
	if err = createInvalidKindName().Execute(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Expected Error! Command Has Invalid(Has Space) API Version Argument.")
		silencelog()
	}

}

func TestValidCommandArgument(t *testing.T) {
	defer os.RemoveAll("./test-operator")
	if err := createValidCommand().Execute(); err != nil {
		resetlog()
		t.Logf("Error: %v", err)
		os.RemoveAll("./test-operator")
		t.Fatal("Unexpected Error! Provided Correct Arguments and Flags.")
		silencelog()
	}

}

// clean up

func TestHelmChartDownload(t *testing.T) {
	var client HelmChartClient

	gopathDir := os.Getenv("GOPATH")
	pathConfig := pathconfig.NewConfig(filepath.Join(gopathDir, "src/github.com/redhat-nfvpe/helm2go-operator-sdk"))
	client.PathConfig = pathConfig

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Panic: %+v\n", r)
	// 	}
	// }()

	client.HelmChartRef = "auto-deploy-app"
	client.HelmChartRepo = "https://charts.gitlab.io/"
	err := client.LoadChart()
	if err != nil {
		resetlog()
		t.Fatalf("Unexpected Error: %v\n", err)
		silencelog()
	}

	client.HelmChartRef = "not-a-chart"
	client.HelmChartRepo = ""
	err = client.LoadChart()
	if err == nil {
		resetlog()
		t.Fatalf("Expected Error Chart Does Not Exist!")
		silencelog()
	}

	client.HelmChartRepo = "invalidrepo.io"
	err = client.LoadChart()
	if err == nil {
		resetlog()
		t.Fatalf("Expected Error Invalid Repo!")
		silencelog()
	}
}

func TestAPINamingConventionValidation(t *testing.T) {
	testCases := map[string]bool{
		"app.example.com/v1alpha1": true,
		"web.k8s.io/v1":            true,
		"k8s.io":                   false,
		"app,s.k8s.io/v1beta1":     false,
		"apps.example/v1beta1":     false,
	}

	for v, ok := range testCases {
		if m := apiVersionMatchesConvention(v); m != ok {
			t.Fatalf("error validating: %v; expected %v got %v", v, ok, m)
		}
	}
}
