package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gopilot",
	Short: "Gopher is a Go version manager",
	Long:  "Gopher lets you list, install, and switch between Go versions easily.",
}
