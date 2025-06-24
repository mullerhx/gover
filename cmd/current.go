package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active Go version",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user home:", err)
			os.Exit(1)
		}

		symlink := filepath.Join(usr.HomeDir, ".gover", "current")

		target, err := os.Readlink(symlink)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No Go version currently active. Use `gover use <version>`.")
			} else {
				fmt.Println("Failed to read symlink:", err)
			}
			os.Exit(1)
		}

		// Extract version from the symlink target
		versionsDir := filepath.Join(usr.HomeDir, ".gover", "versions") + string(os.PathSeparator)
		version := target
		if strings.HasPrefix(target, versionsDir) {
			version = target[len(versionsDir):]
		}

		fmt.Println("Current Go version:", version)
		fmt.Println("GOROOT:", symlink)
	},
}

func init() {
	RootCmd.AddCommand(currentCmd)
}
