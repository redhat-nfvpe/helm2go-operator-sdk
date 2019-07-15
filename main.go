package main

import (
	"os"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/cmd/helm2go-operator-sdk/convert"
	"github.com/spf13/cobra"
)

func main() {
	// cwd, _ := os.Getwd()
	// resourcesPath := filepath.Join(cwd, "test", "resources")
	// rs, err := load.YAMLUnmarshalResources(resourcesPath)
	// if err != nil {
	// 	panic(err)
	// }
	// err = templating.TestDeclarationTemplate(rs[0])
	// if err != nil {
	// 	fmt.Errorf("%v", err)
	// }
	// // litter.Dump(rs[0])
	// s := fmt.Sprintf("%#+v", rs[0])
	// fmt.Println(s)

	// this is the code that i actually need to use
	root := &cobra.Command{
		Use:   "helm2go-operator-sdk",
		Short: "A Kit to Convert Helm Chart Operators to Go Operators, Yee-Haw! 🏇",
	}
	root.AddCommand(convert.NewConvertCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
