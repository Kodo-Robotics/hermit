/*
Copyright © 2025 Kodo Robotics
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
	"github.com/spf13/cobra"
)

var outputBoxPath string

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the current VM into a .box file",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		vmName := cfg.Name
		outputDir := "hermit_build"
		os.MkdirAll(outputDir, 0755)

		fmt.Println("📦 Exporting VM...")
		ovfPath, err := virtualbox.ExportVM(vmName, outputDir)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		diskPath, err := virtualbox.FindDiskFile(outputDir)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		boxName := outputBoxPath
		if boxName == "" {
			boxName = vmName + ".box"
		}
		fmt.Println("📦 Creating", boxName)

		files := map[string]string{
			"box.ovf":               ovfPath,
			filepath.Base(diskPath): diskPath,
		}

		if err := virtualbox.CreateBoxArchive(boxName, files); err != nil {
			fmt.Println("❌", err)
			return
		}

		fmt.Println("🧹 Cleaning up temporary files...")
		if err := os.RemoveAll(outputDir); err != nil {
			fmt.Println("⚠️ Warning: failed to clean temp directory:", err)
		}

		fmt.Println("✅ VM packaged into", boxName)
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.Flags().StringVarP(&outputBoxPath, "output", "o", "", "Path to output .box file")
}
