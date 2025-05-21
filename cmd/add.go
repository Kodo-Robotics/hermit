/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/utils"
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
		destDir := filepath.Join(".hermit", "boxes", boxName)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Println("❌ Failed to create box directory:", err)
			return
		}

		fmt.Println("📦 Extracting box to:", destDir)
		if err := utils.ExtractTar(boxFile, destDir); err != nil {
			fmt.Println("❌ Failed to extract .box:", err)
			return
		}

		fmt.Printf("✅ Box '%s' added successfully.\n", boxName)
	},
}

func init() {
	boxCmd.AddCommand(addCmd)
}