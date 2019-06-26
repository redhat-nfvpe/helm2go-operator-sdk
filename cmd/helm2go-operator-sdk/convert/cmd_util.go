package convert

import (
	"fmt"

	"github.com/spf13/cobra"
)

func verifyFlags() error {
	// TODO PLACEHOLDER
	return nil
}

func parse(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Please Use Specified Flags")
	}
	return nil
}
