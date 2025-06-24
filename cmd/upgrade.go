package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade <major-version>",
	Short: "Upgrade to the latest patch release of a major Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		major := args[0]
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user info:", err)
			os.Exit(1)
		}

		// Load cached releases
		cachePath := filepath.Join(usr.HomeDir, ".gover", "releases.json")
		file, err := os.Open(cachePath)
		if err != nil {
			fmt.Println("Failed to open releases cache:", err)
			os.Exit(1)
		}
		defer file.Close()

		var versions []GoVersion
		if err := json.NewDecoder(file).Decode(&versions); err != nil {
			fmt.Println("Failed to decode releases cache:", err)
			os.Exit(1)
		}

		// Find all stable versions matching the major prefix
		prefix := "go" + major + "."
		var matches []string
		for _, v := range versions {
			if v.Stable && strings.HasPrefix(v.Version, prefix) {
				matches = append(matches, v.Version)
			}
		}

		if len(matches) == 0 {
			fmt.Printf("No stable versions found for major version %s\n", major)
			os.Exit(1)
		}

		sort.Strings(matches)
		latest := matches[len(matches)-1]

		// Check if installed
		installPath := filepath.Join(usr.HomeDir, ".gover", "versions", latest)
		if !fileExists(installPath) {
			fmt.Printf("Version %s not installed. Installing...\n", latest)
			// Call your existing install logic here, e.g.:
			err := installVersion(latest)
			if err != nil {
				fmt.Println("Installation failed:", err)
				os.Exit(1)
			}
			fmt.Println("Installation complete.")
		} else {
			fmt.Printf("Version %s is already installed.\n", latest)
		}

		// Switch to latest
		fmt.Printf("Switching to %s...\n", latest)
		err = switchVersion(latest)
		if err != nil {
			fmt.Println("Failed to switch version:", err)
			os.Exit(1)
		}
		fmt.Println("Upgrade complete.")
	},
}

func installVersion(version string) error {
	usr, _ := user.Current()
	destDir := filepath.Join(usr.HomeDir, ".gover", "versions", version)

	if fileExists(destDir) {
		return nil
	}

	url := fmt.Sprintf("https://go.dev/dl/%s.%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Download to temp
	tmpFile, err := os.CreateTemp("", version+".*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save archive: %w", err)
	}
	tmpFile.Close()

	// Extract
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("tar", "-C", destDir, "--strip-components=1", "-xzf", tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	return nil
}

func switchVersion(version string) error {
	usr, _ := user.Current()
	targetPath := filepath.Join(usr.HomeDir, ".gover", "versions", version)
	symlinkPath := filepath.Join(usr.HomeDir, ".gover", "current")

	if !fileExists(targetPath) {
		return fmt.Errorf("version not installed: %s", version)
	}

	_ = os.Remove(symlinkPath) // Remove existing symlink if any

	err := os.Symlink(targetPath, symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Output shell eval if in use mode
	fmt.Printf(`export GOROOT="%s"
export PATH="%s/bin:$PATH"
export GOPATH="%s/go"
`, symlinkPath, symlinkPath, usr.HomeDir)

	return nil
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
}
