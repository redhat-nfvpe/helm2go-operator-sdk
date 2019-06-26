package convert

import (
	"os"
	"path/filepath"
	"testing"
)

var cwd, _ = os.Getwd()
var parent = filepath.Dir(filepath.Dir(filepath.Dir(cwd)))
var local = "/test/bitcoind"
var testLocal = parent + local

func TestLoadChartGetsChart(t *testing.T) {
	// point test to right directory
	helmChartRef = testLocal
	// load the chart
	loadChart()
	// verify that the chart loads the right thing
	if chartName != "bitcoind" {
		t.Fatalf("Unexpected Chart Name!")
	}
}
