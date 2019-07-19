package main

import (
	"os"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/cmd/helm2go-operator-sdk/convert"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "helm2go-operator-sdk",
		Short: "A Kit to Convert Helm Chart Operators to Go Operators, Yee-Haw! 🏇",
	}
	root.AddCommand(convert.NewConvertCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
