package convert

import (
	"fmt"
	"strings"
)

func verifyFlags() error {
	if len(helmChartRef) == 0 {
		return fmt.Errorf("Please Specify Helm Chart Reference")
	}
	if strings.ContainsAny(helmChartRef, " ") {
		return fmt.Errorf("Helm Chart Cannot Contain Spaces")
	}
	if len(kind) == 0 {
		return fmt.Errorf("Please Specify Operator Kind")
	}
	if strings.ContainsAny(kind, " ") {
		return fmt.Errorf("Kind Cannot Contain Spaces")
	}
	if len(apiVersion) == 0 {
		return fmt.Errorf("Please Specify Operator API Version")
	}
	if strings.ContainsAny(apiVersion, " ") {
		return fmt.Errorf("API Version Cannot Contain Spaces")
	}
	return nil
}

func parse(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Please Specify Operator Name")
	}
	outputDir = args[0]
	if len(outputDir) == 0 {
		return fmt.Errorf("Project Name Must Not Be Empty")
	}
	return nil
}
