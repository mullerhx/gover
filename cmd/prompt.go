package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Output current Go version for shell prompt",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			os.Exit(0) // Silent fail for prompt integration
		}
		current := filepath.Join(usr.HomeDir, ".gover", "current")
		resolved, err := filepath.EvalSymlinks(current)
		if err != nil {
			os.Exit(0)
		}
		version := filepath.Base(resolved)
		if strings.HasPrefix(version, "go") {
			fmt.Printf("[🐹 %s] ", version)
		}
	},
}

func init() {
	RootCmd.AddCommand(promptCmd)
}
