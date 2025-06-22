package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Download and install a specific Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		url := fmt.Sprintf("https:/golang.org/dl/%s.%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
		fmt.Println("Downloading:", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error downloading:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to download: %s\n", resp.Status)
			os.Exit(1)
		}

		outFile := filepath.Join(os.TempDir(), version+".tar.gz")
		out, err := os.Create(outFile)
		if err != nil {
			fmt.Println("Failed to create temp file:", err)
			os.Exit(1)
		}
		defer out.Close()

		total := resp.ContentLength
		progressReader := &progressReader{Reader: resp.Body, total: total}
		if _, err := io.Copy(out, progressReader); err != nil {
			fmt.Println("Download failed:", err)
			os.Exit(1)
		}

		fmt.Println("\nExtracting...")
		if err := extractTarGz(outFile, filepath.Join(os.Getenv("HOME"), ".gopilot", "versions", version)); err != nil {
			fmt.Println("Failed to extract archive:", err)
			os.Exit(1)
		}

		fmt.Println("Installation completed successfully.")
	},
}

type progressReader struct {
	Reader io.Reader
	total  int64
	read   int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.read += int64(n)
	fmt.Printf("\rProgress: %.2f%%", float64(pr.read)/float64(pr.total)*100)
	return n, err
}

func extractTarGz(file, targetDir string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(targetDir, hdr.Name)
		if strings.HasSuffix(hdr.Name, "/") {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		outFile, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
	}
	return nil
}

func init() {
	RootCmd.AddCommand(installCmd)
}
