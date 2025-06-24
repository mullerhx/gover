package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var autoUse = false

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if autoUse {
			version, err := detectGoModVersion()
			if err != nil {
				fmt.Println("Auto detection failed:", err)
				os.Exit(1)
			}
			args = []string{version}
		}

		if len(args) < 1 {
			fmt.Println("Usage: gover use <version> or --auto")
			os.Exit(1)
		}
		version := args[0]
		if !strings.HasPrefix(version, "go") {
			version = "go" + version
		}

		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get current user:", err)
			os.Exit(1)
		}

		installPath := filepath.Join(usr.HomeDir, ".gover", "versions", version, "go")
		if _, err := os.Stat(installPath); os.IsNotExist(err) {
			fmt.Printf("Version %s not installed. Run `gover install %s` first.\n", version, version)
			os.Exit(1)
		}

		currentLink := filepath.Join(usr.HomeDir, ".gover", "current")
		_ = os.Remove(currentLink)
		if err := os.Symlink(installPath, currentLink); err != nil {
			fmt.Println("Failed to create symlink:", err)
			os.Exit(1)
		}
		currentBinLink := filepath.Join(usr.HomeDir, ".gover", "current", "bin")
		files, err := os.ReadDir(currentBinLink)
		if err != nil {
			fmt.Println("Failed to read bin directory:", err)
		} else {
			for _, file := range files {
				if !file.IsDir() {
					filePath := filepath.Join(currentBinLink, file.Name())
					err := os.Chmod(filePath, 0755)
					if err != nil {
						fmt.Printf("Failed to chmod %s: %v\n", filePath, err)
					}
				}
			}
		}

		toolsLink := filepath.Join(usr.HomeDir, ".gover", "current", "pkg", "tool", runtime.GOOS+"_"+runtime.GOARCH)
		files, err = os.ReadDir(toolsLink)
		if err != nil {
			fmt.Println("Failed to read ", toolsLink, " directory:", err)
		} else {
			for _, file := range files {
				if !file.IsDir() {
					filePath := filepath.Join(toolsLink, file.Name())
					err := os.Chmod(filePath, 0755)
					if err != nil {
						fmt.Printf("Failed to chmod %s: %v\n", filePath, err)
					}
				}
			}
		}

		shell := detectShell()
		profilePath := shellProfile(shell)

		fmt.Println("✅ Go version", version, "is now active via ~/.gover/current")
		fmt.Println("👉 Add the following to your", profilePath, "if not already present:\n")

		fmt.Println("export GOROOT=\"$HOME/.gover/current\"")
		fmt.Println("export PATH=\"$HOME/.gover/current/bin:$PATH\"")
		fmt.Println("export GOPATH=\"$HOME/go\"")
	},
}

func init() {
	useCmd.Flags().BoolVar(&autoUse, "auto", false, "Automatically detect version from go.mod")
	RootCmd.AddCommand(useCmd)
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "bash") {
		return "bash"
	}
	if strings.Contains(shell, "fish") {
		return "fish"
	}
	return "sh"
}

func shellProfile(shell string) string {
	usr, _ := user.Current()
	switch shell {
	case "zsh":
		return filepath.Join(usr.HomeDir, ".zshrc")
	case "bash":
		return filepath.Join(usr.HomeDir, ".bashrc")
	case "fish":
		return filepath.Join(usr.HomeDir, ".config/fish/config.fish")
	default:
		return filepath.Join(usr.HomeDir, ".profile")
	}
}

func detectGoModVersion() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if fileExists(goModPath) {
			data, err := os.ReadFile(goModPath)
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "go ") {
					version := strings.TrimPrefix(line, "go ")
					version = strings.TrimSpace(version)
					return resolveLatestPatch("go" + version)
				}
			}
			break
		}

		// move one directory up
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("no go.mod found")
}

func resolveLatestPatch(prefix string) (string, error) {
	usr, _ := user.Current()
	cachePath := filepath.Join(usr.HomeDir, ".gover", "releases.json")

	var versions []GoVersion
	file, err := os.Open(cachePath)
	if err != nil {
		return "", fmt.Errorf("failed to read release cache: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&versions); err != nil {
		return "", err
	}

	var matches []string
	for _, v := range versions {
		if strings.HasPrefix(v.Version, prefix) && v.Stable {
			matches = append(matches, v.Version)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching versions for prefix %s", prefix)
	}

	sort.Strings(matches)
	return matches[len(matches)-1], nil
}
