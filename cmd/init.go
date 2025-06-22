package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		OS       string `json:"os"`
		Arch     string `json:"arch"`
	} `json:"files"`
}

var all bool
var majorFilter string
var forceFetch bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user info:", err)
			os.Exit(1)
		}

		releasesPath := filepath.Join(usr.HomeDir, ".gopilot", "releases.json")
		var versions []GoVersion

		if forceFetch || !fileExists(releasesPath) {
			resp, err := http.Get("https://golang.org/dl/?mode=json")
			if err != nil {
				fmt.Println("Failed to fetch versions:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
				fmt.Println("Failed to decode response:", err)
				os.Exit(1)
			}

			if err := os.MkdirAll(filepath.Dir(releasesPath), 0755); err != nil {
				fmt.Println("Failed to create .gopilot directory:", err)
				os.Exit(1)
			}

			file, err := os.Create(releasesPath)
			if err != nil {
				fmt.Println("Failed to create releases file:", err)
				os.Exit(1)
			}
			defer file.Close()
			_ = json.NewEncoder(file).Encode(versions)
		} else {
			file, err := os.Open(releasesPath)
			if err != nil {
				fmt.Println("Failed to read local releases file:", err)
				os.Exit(1)
			}
			defer file.Close()
			if err := json.NewDecoder(file).Decode(&versions); err != nil {
				fmt.Println("Failed to decode local releases file:", err)
				os.Exit(1)
			}
		}

		versionMap := map[string][]string{}

		for _, v := range versions {
			if !all && !v.Stable {
				continue
			}
			for _, f := range v.Files {
				if f.OS == runtime.GOOS && f.Arch == runtime.GOARCH {
					parts := strings.Split(v.Version, ".")
					if len(parts) < 2 {
						continue
					}
					prefix := strings.Join(parts[:2], ".")
					cleanPrefix := strings.TrimPrefix(prefix, "go")
					if majorFilter != "" && !strings.HasPrefix(cleanPrefix, majorFilter) {
						continue
					}
					versionMap[prefix] = append(versionMap[prefix], v.Version)
				}
			}
		}

		var keys []string
		for k := range versionMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, prefix := range keys {
			vers := versionMap[prefix]
			sort.Strings(vers)
			limit := 4
			if len(vers) < 4 {
				limit = len(vers)
			}
			for _, v := range vers[len(vers)-limit:] {
				fmt.Println(v)
			}
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gopilot environment",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user info:", err)
			os.Exit(1)
		}
		dir := filepath.Join(usr.HomeDir, ".gopilot")
		releasesPath := filepath.Join(dir, "releases.json")

		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Failed to create .gopilot directory:", err)
			os.Exit(1)
		}

		fmt.Println("Fetching release list...")
		resp, err := http.Get("https://golang.org/dl/?mode=json")
		if err != nil {
			fmt.Println("Failed to fetch versions:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var versions []GoVersion
		if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
			fmt.Println("Failed to decode response:", err)
			os.Exit(1)
		}

		file, err := os.Create(releasesPath)
		if err != nil {
			fmt.Println("Failed to create releases file:", err)
			os.Exit(1)
		}
		defer file.Close()
		if err := json.NewEncoder(file).Encode(versions); err != nil {
			fmt.Println("Failed to write releases file:", err)
			os.Exit(1)
		}
		fmt.Println("Gopilot initialized successfully.")
	},
}
