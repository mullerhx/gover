package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gover",
	Short: "Gopliot is a Go version manager",
	Long:  "Gopliot lets you list, install, and switch between Go versions easily.",
}
