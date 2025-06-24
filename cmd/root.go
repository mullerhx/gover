package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gover",
	Short: "Gover is a Go version manager",
	Long:  "Gover lets you list, install, and switch between Go versions easily.",
}
