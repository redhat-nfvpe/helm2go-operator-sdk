package main

import (
	"fmt"
	"reflect"

	"github.com/redhat-nvfpe/helm2go-operator-sdk/pkg/convert"
)

func main() {
	convert.DirectoryInjectedYAMLToJSON("./test/resources")
	testDep, err := convert.JSONUnmarshal("./test/jsonOutputs/test_Deployment_0.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(reflect.TypeOf(testDep).String())
}
