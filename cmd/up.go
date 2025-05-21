/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/utils"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Boot up the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		boxDir := filepath.Join(".hermit", "boxes", strings.ReplaceAll(cfg.Box, "/", "_"))
		diskPath, err := utils.FindDiskImage(boxDir)
		if err != nil {
			fmt.Println("❌", err)
			fmt.Println("👉 Run `hermit box add <path>.box` to install the box.")
			return
		}

		fmt.Println("🚀 Launching VM...")
		err = virtualbox.CreateAndStartVM(
			cfg.Name, 
			cfg.Memory, 
			cfg.CPUs, 
			cfg.VRAM,
			cfg.GraphicsController,
			diskPath,
		)
		if err != nil {
			fmt.Println("❌ Error starting VM:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}