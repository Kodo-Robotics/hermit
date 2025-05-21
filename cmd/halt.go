/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
)

var haltCmd = &cobra.Command{
	Use:   "halt",
	Short: "Gracefully shut down the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()

		state, err := virtualbox.GetVMState(cfg.Name)
		if err != nil {
			fmt.Println("❌ Could not determine VM state:", err)
			return
		}

		if state == "poweroff" {
			fmt.Println("⏹️ VM is already stopped.")
			return
		}

		fmt.Println("🛑 Sending shutdown signal to VM...")
		err = virtualbox.HaltVM(cfg.Name)
		if err != nil {
			fmt.Printf("❌ Failed to halt VM: %v\n", err)
			return
		}

		fmt.Println("✅ Shutdown signal sent. VM will power off shortly.")
	},
}

func init() {
	rootCmd.AddCommand(haltCmd)
}
