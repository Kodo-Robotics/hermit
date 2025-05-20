/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/config"
)

var initCmd = &cobra.Command{
	Use:   "init [box-name]",
	Short: "Initialize a new Hermit VM configuration file (hermit.json)",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("hermit.json"); err == nil {
			fmt.Println("⚠️  hermit.json already exists. Use --force to overwrite (not yet supported).")
			return
		}

		boxName := "ubuntu/focal64"
		if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
			boxName = args[0]
		}

		cfg := config.HermitConfig{
			Box:			boxName,
			Provider:		"virtualbox",
			Name:			"hermit-vm",
			CPUs:			1,
			Memory:			1024,
			DiskSizeMB:		10000,
			ForwardedPorts:	[]config.Port {
				{Guest: 22, Host: 2222},
			},
		}

		jsonData, err := json.MarshalIndent(cfg, "", " ")
		if err != nil {
			fmt.Println("❌ Could not write hermit.json:", err)
			return
		}

		err = os.WriteFile("hermit.json", jsonData, 0644)
		if err != nil {
			fmt.Println("❌ Could not write hermit.json:", err)
			return
		}

		absPath, _ := filepath.Abs("hermit.json")
		fmt.Println("🌱 Created hermit.json at", absPath)
		fmt.Println("📦 Next: Run `hermit up` to start the VM")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
