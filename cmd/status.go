/*
Copyright © 2025 Kodo Robotics
*/
package cmd

import (
	"fmt"

	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		state, err := virtualbox.GetVMState(cfg.Name)
		if err != nil {
			fmt.Printf("❌ VM '%s' not found in VirtualBox.\n", cfg.Name)
			return
		}

		fmt.Printf("🖥️ VM Name: %s\n", cfg.Name)
		fmt.Printf("📦 Box:     %s\n", cfg.Box)
		fmt.Printf("⚙️  State:   %s\n", state)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
