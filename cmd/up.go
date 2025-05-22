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

		// Check if VM already exists
		_, err = virtualbox.GetVMState(cfg.Name)
		if err == nil {
			if err := virtualbox.StartVM(cfg.Name); err != nil {
				fmt.Println("❌ Failed to start VM:", err)
			}
			return
		} else {
			// VM not found - check for stale vbox
			if err := virtualbox.CleanupStaleVBoxFile(cfg.Name); err != nil {
				fmt.Println("⚠️ Warning: could not clean stale .vbox file:", err)
			}
		}
		

		boxDir := filepath.Join(".hermit", "boxes", strings.ReplaceAll(cfg.Box, "/", "_"))
		ovfPath, err := utils.FindOVF(boxDir)
		if err != nil {
			fmt.Println("❌", err)
			fmt.Println("👉 Run `hermit box add <path>.box` to install the box.")
			return
		}

		fmt.Printf("📦 Importing VM '%s' from box: %s\n", cfg.Name, cfg.Box)
		if err := virtualbox.ImportOVF(ovfPath, cfg.Name); err != nil {
			fmt.Println("❌ Failed to import OVF:", err)
			return
		}

		fmt.Println("⚙️ Applying CPU and memory settings...")
		if err := virtualbox.ModifyVM(cfg.Name, cfg.Memory, cfg.CPUs); err != nil {
			fmt.Println("⚠️ Failed to modify VM settings:", err)
		}

		for _, port := range cfg.ForwardedPorts {
			fmt.Printf("🔁 Forwarding host:%d -> guest:%d\n", port.Host, port.Guest)
			if err := virtualbox.AddPortForward(cfg.Name, port.Guest, port.Host); err != nil {
				fmt.Printf("⚠️ Failed to add port forward: %v\n", err)
			}
		}

		fmt.Println("🚀 Starting VM...")
		if err := virtualbox.StartVM(cfg.Name); err != nil {
			fmt.Println("❌ Error starting VM:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}