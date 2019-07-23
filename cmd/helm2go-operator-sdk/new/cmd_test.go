package new

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

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

func TestLoadChartGetsLocalChart(t *testing.T) {
	//g := gomega.NewGomegaWithT(t)
	silencelog()

	var local = "/test/bitcoind"
	var testLocal = parent + local

	// load the chart
	chartClient := NewChartClient()
	chartClient.HelmChartRef = testLocal
	err := chartClient.LoadChart()
	if err != nil {
		t.Fatal(err)
	}
	// verify that the chart loads the right thing
	if chartClient.Chart.GetMetadata().Name != "bitcoind" {
		resetlog()
		t.Fatalf("Unexpected Chart Name!")
	}
}

func TestCommandFlagValidation(t *testing.T) {
	// initiate flags for testing
	var err error
	helmChartRef = ""
	silencelog()
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Expected Flag Validation Error: --helm-chart-ref")
		silencelog()
	}
	helmChartRef = "./test/tomcat"
	kind = "Tomcat"
	apiVersion = ""
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Expected Flag Validation Error: --api-version")
		silencelog()
	}
	apiVersion = "app.example/v1alpha1"
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("expected flag validation error; api version does not match naming convention")
		silencelog()
	}
	apiVersion = "app.example.com/v1alpha1"
	if err = verifyFlags(); err != nil {
		resetlog()
		t.Logf("Error: %v", err)
		t.Fatal("Unexpected Flag Validation Error")
		silencelog()
	}
}

func TestKindTypeValidation(t *testing.T) {
	var err error
	apiVersion = "app.example.com/v1alpha1"
	helmChartRef = "./test/tomcat"
	kind = ""
	silencelog()
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Fatal("Expected Flag Validation Error: kind")
		silencelog()
	}
	kind = "tensorflow-notebook"
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Fatalf("Expected Lowercase Flag Validation Error: --kind %v", err)
		silencelog()
	}
	kind = "tensorflowNotebook"
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Fatal("Expected Lowercase Flag Validation Error: --kind")
		t.Fatal(err)
		silencelog()
	}
	kind = "Tensorflow-Notebook"
	if err = verifyFlags(); err == nil {
		resetlog()
		t.Fatal("Expected Hyphen Flag Validation Error: --kind")
		t.Fatal(err)
		silencelog()
	}
	kind = "TensorflowNotebook"
	if err = verifyFlags(); err != nil {
		resetlog()
		t.Fatal("Unexpected Flag Validation Error")
		t.Fatal(err)
		silencelog()
	}
}

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
