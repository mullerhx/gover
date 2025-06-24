package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect Go version from nearest go.mod and resolve latest patch version",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := detectGoModVersion()
		if err != nil {
			fmt.Println("❌", err)
			os.Exit(1)
		}
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(detectCmd)
}
