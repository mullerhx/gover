package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
)

var force bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user info:", err)
			os.Exit(1)
		}

		installPath := filepath.Join(usr.HomeDir, ".gover", "versions", version)

		currentPath, _ := filepath.EvalSymlinks(filepath.Join(usr.HomeDir, ".gover", "current"))

		if currentPath == installPath && !force {
			fmt.Printf("⚠️  %s is currently in use. Use --force to uninstall it anyway.\n", version)
			os.Exit(1)
		}

		if !fileExists(installPath) {
			fmt.Printf("❌ Version %s is not installed.\n", version)
			os.Exit(1)
		}

		err = os.RemoveAll(installPath)
		if err != nil {
			fmt.Println("Failed to uninstall:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Uninstalled %s\n", version)
	},
}

func init() {
	uninstallCmd.Flags().BoolVarP(&force, "force", "", false, "Force uninstall even if version is active")
	RootCmd.AddCommand(uninstallCmd)
}
