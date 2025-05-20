/*
Copyright © 2025 Kodo-Robotics

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type HermitConfig struct {
	Name		string	`json:"name"`
	CPUs		int		`json:"cpus"`
	Memory		int		`json:"memory"`
	DiskSizeMB	int		`json:"disk_size_mb"`
	ISOPath		string	`json:"iso_path"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Hermit VM configuration file (hermit.json)",
	Run: func(cmd *cobra.Command, args []string) {
		config := HermitConfig{
			Name:		"hermit-vm",
			CPUs:		1,
			Memory:		1024,
			DiskSizeMB:	10000,
			ISOPath:	"/absolute/path/to/ubuntu.iso",
		}

		// Check if file already exists
		if _, err := os.Stat("hermit.json"); err == nil {
			fmt.Println("⚠️  hermit.json already exists. Aborting.")
			return
		}

		jsonData, err := json.MarshalIndent(config, "", " ")
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
		fmt.Println("📦 Next: Edit the ISO path, then run `hermit up`")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
