/*
Copyright © 2025 Kodo Robotics
*/
package boxcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kodo-Robotics/hermit/pkg/core"
	"github.com/Kodo-Robotics/hermit/pkg/utils"
	"github.com/spf13/cobra"
)

var boxAlias string

var addCmd = &cobra.Command{
	Use:   "add <box-file>",
	Short: "Add a Hermit box from a .box file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ Please provide a path to the .box file")
			return
		}

		boxFile := args[0]

		if _, err := os.Stat(boxFile); os.IsNotExist(err) {
			fmt.Println("❌ Box file not found:", boxFile)
			return
		}

		boxName := boxAlias
		if boxName == "" {
			boxName = strings.TrimSuffix(filepath.Base(boxFile), filepath.Ext(boxFile))
		}
		destDir := utils.GetBoxPath(boxName)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			fmt.Println("❌ Failed to create box directory:", err)
			return
		}

		fmt.Println("📦 Extracting box to:", destDir)
		if err := utils.ExtractTar(boxFile, destDir); err != nil {
			fmt.Println("❌ Failed to extract .box:", err)
			return
		}

		// Register box
		if err := core.AddBox(boxName, destDir); err != nil {
			fmt.Println("❌ Failed to register box in registry:", err)
			return
		}

		fmt.Printf("✅ Box '%s' added successfully.\n", boxName)
	},
}

func init() {
	addCmd.Flags().StringVarP(&boxAlias, "name", "n", "", "Custom name (alias) for this box")
	BoxCmd.AddCommand(addCmd)
}
