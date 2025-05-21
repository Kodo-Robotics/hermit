/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/config"
)

var haltCmd = &cobra.Command{
	Use:   "halt",
	Short: "Gracefully shut down the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()

		fmt.Println("🛑 Sending shutdown signal to VM...")
		cmdShutdown := exec.Command("VBoxManage", "controlvm", cfg.Name, "acpipowerbutton")
		cmdShutdown.Stdout = os.Stdout
		cmdShutdown.Stderr = os.Stderr
		err = cmdShutdown.Run()

		if err != nil {
			fmt.Printf("❌ Failed to send shutdown signal: %v\n", err)
			return
		}

		fmt.Println("✅ Shutdown signal sent. VM will power off shortly.")
	},
}

func init() {
	rootCmd.AddCommand(haltCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// haltCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// haltCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
