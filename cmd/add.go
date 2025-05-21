/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <box-file>",
	Short: "Add a Hermit box from a .box file",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		boxFile := args[0]

		if _, err := os.Stat(boxFile); os.IsNotExist(err) {
			fmt.Println("❌ Box file not found:", boxFile)
			return
		}

		baseName := strings.TrimSuffix(filepath.Base(boxFile), filepath.Ext(boxFile))
		boxName := strings.ReplaceAll(baseName, ".", "_")
		destDir := filepath.join(".hermit", "boxes", boxName)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Println("❌ Failed to create box directory:", err)
			return
		}

		fmt.Println("📦 Extracting box to:", destDir)
		if err := extractTar(boxFile, destDir); err != nil {
			fmt.Println("❌ Failed to extract .box:", err)
			return
		}

		fmt.Printf("✅ Box '%s' added successfully.\n", boxName)
	},
}

func init() {
	box.AddCommand(addCmd)
}

func extractTar(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tgz") {
		gz, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gz.Close()
		reader = gz
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(targetPath, 0755)
		case tar.TypeReg:
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			fmt.Printf("⚠️ Skipping unsupported file: %s\n", header.Name)
		}
	}

	return nil
}
