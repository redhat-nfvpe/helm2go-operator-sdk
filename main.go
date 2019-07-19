package main

import (
	"os"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/cmd/helm2go-operator-sdk/newcmd"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "helm2go-operator-sdk",
		Short: "A Kit to Create Go Operators for Helm Charts, Yee-Haw! 🏇",
	}
	root.AddCommand(newcmd.NewCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
