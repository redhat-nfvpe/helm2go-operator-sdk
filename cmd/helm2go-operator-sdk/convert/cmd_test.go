package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

var cwd, _ = os.Getwd()
var parent = filepath.Dir(filepath.Dir(filepath.Dir(cwd)))

func TestLoadChartGetsLocalChart(t *testing.T) {
	var local = "/test/bitcoind"
	var testLocal = parent + local

	// point test to right directory
	helmChartRef = testLocal
	// load the chart
	loadChart()
	// verify that the chart loads the right thing
	if chartName != "bitcoind" {
		t.Fatalf("Unexpected Chart Name!")
	}
}

func TestCommandFlagValidation(t *testing.T) {
	// initiate flags for testing
	var err error
	helmChartRef = ""
	if err = verifyFlags(); err == nil {
		t.Logf("Error: %v", err)
		t.Fatal("Expected Flag Validation Error: --helm-chart-ref")
	}
	helmChartRef = "./test/tomcat"
	kind = "Tomcat"
	apiVersion = ""
	if err = verifyFlags(); err == nil {
		t.Logf("Error: %v", err)
		t.Fatal("Expected Flag Validation Error: --api-version")
	}
	apiVersion = "app.example.com/v1alpha1"
	if err = verifyFlags(); err != nil {
		t.Logf("Error: %v", err)
		t.Fatal("Unexpected Flag Validation Error")
	}
}

func TestKindTypeValidation(t *testing.T) {
	var err error
	apiVersion = "app.example.com/v1alpha1"
	helmChartRef = "./test/tomcat"
	kind = ""
	if err = verifyFlags(); err == nil {
		t.Logf("Error: %v", err)
		t.Fatal("Expected Flag Validation Error: kind")
	}
	kind = "tensorflow-notebook"
	if err = verifyFlags(); err == nil {
		t.Fatal("Expected Lowercase Flag Validation Error: --kind")
	}
	kind = "tensorflowNotebook"
	if err = verifyFlags(); err == nil {
		t.Fatal("Expected Lowercase Flag Validation Error: --kind")
	}
	kind = "Tensorflow-Notebook"
	if err = verifyFlags(); err == nil {
		t.Fatal("Expected Hyphen Flag Validation Error: --kind")
	}
	kind = "TensorflowNotebook"
	if err = verifyFlags(); err != nil {
		t.Fatal("Unexpected Flag Validation Error")
	}
}

//creates a command for use in the argument validation test
func createValidCommand() *cobra.Command {
	var cmd *cobra.Command
	cmd = NewConvertCmd()

	tomcatDir := "/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/tomcat"

	// command has an operator name and correct flags
	cmd.SetArgs([]string{
		"test-operator",
		fmt.Sprintf("--helm-chart=%s", tomcatDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Tomcat"),
	})
	fmt.Println(cmd)
	return cmd
}

// creates an invalid command for use in the argument validation test
func createInvalidCommand() *cobra.Command {
	var cmd *cobra.Command
	cmd = NewConvertCmd()

	tomcatDir := "/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/tomcat"

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		fmt.Sprintf("--helm-chart=%s", tomcatDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Tomcat"),
	})
	return cmd
}

func createInvalidOperatorName() *cobra.Command {
	var cmd *cobra.Command
	cmd = NewConvertCmd()

	tomcatDir := "/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/tomcat"

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		"",
		fmt.Sprintf("--helm-chart=%s", tomcatDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1"),
		fmt.Sprintf("--kind=%s", "Tomcat"),
	})
	return cmd
}

func createInvalidKindName() *cobra.Command {
	var cmd *cobra.Command
	cmd = NewConvertCmd()

	tomcatDir := "/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/tomcat"

	// command does not have the operator name; invalid
	cmd.SetArgs([]string{
		"test-operator",
		fmt.Sprintf("--helm-chart=%s", tomcatDir),
		fmt.Sprintf("--api-version=%s", "app.example.com/v1alpha1 Server"),
		fmt.Sprintf("--kind=%s", "Tomcat"),
	})
	return cmd
}

func TestCommandArgumentValidation(t *testing.T) {
	var err error
	if err = createInvalidCommand().Execute(); err == nil {
		t.Logf("Error: %v", err)
		os.RemoveAll("./test-operator")
		t.Fatal("Expected Error! Command Has No Operator Name Argument.")
	}
	if err = createInvalidOperatorName().Execute(); err == nil {
		t.Logf("Error: %v", err)
		os.RemoveAll("./test-operator")
		t.Fatal("Expected Error! Command Has Invalid(empty) Operator Name Argument.")
	}
	if err = createInvalidKindName().Execute(); err == nil {
		t.Logf("Error: %v", err)
		os.RemoveAll("./test-operator")
		t.Fatal("Expected Error! Command Has Invalid(Has Space) API Version Argument.")
	}
	if err = createValidCommand().Execute(); err != nil {
		t.Logf("Error: %v", err)
		os.RemoveAll("./test-operator")
		t.Fatal("Unexpected Error! Provided Correct Arguments and Flags.")
	}

	// clean up
	os.RemoveAll("./test-operator")
}

func TestHelmChartDownload(t *testing.T) {
	helmChartRef = "auto-deploy-app"
	helmChartRepo = "https://charts.gitlab.io/"
	outputDir = "test-operator"
	err := loadChart()
	if err != nil {
		t.Fatalf("Unexpected Error: %v\n", err)
	}

	helmChartRef = "not-a-chart"
	err = loadChart()
	if err == nil {
		t.Fatalf("Expected Error Chart Does Not Exist!")
	}

	helmChartRepo = "invalidrepo.io"
	err = loadChart()
	if err == nil {
		t.Fatalf("Expected Error Invalid Repo!")
	}
}
