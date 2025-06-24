package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gover environment",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get user info:", err)
			os.Exit(1)
		}
		dir := filepath.Join(usr.HomeDir, ".gover")
		releasesPath := filepath.Join(dir, "releases.json")

		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Failed to create .gover directory:", err)
			os.Exit(1)
		}

		fmt.Println("Fetching release list...")
		resp, err := http.Get("https://golang.org/dl/?mode=json&include=all")
		if err != nil {
			fmt.Println("Failed to fetch versions:", err)
			os.Exit(1)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

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
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		if err := json.NewEncoder(file).Encode(versions); err != nil {
			fmt.Println("Failed to write releases file:", err)
			os.Exit(1)
		}
		fmt.Println("Gover initialized successfully.")
	},
}
