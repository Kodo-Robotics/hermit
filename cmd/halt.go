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

var haltCmd = &cobra.Command{
	Use:   "halt",
	Short: "Gracefully shut down the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		err = virtualbox.HaltVM(cfg.Name)
		if err != nil {
			fmt.Println("❌", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(haltCmd)
}
