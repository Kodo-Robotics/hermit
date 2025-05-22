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

		err = virtualbox.HaltVM(cfg.Name)
		if err != nil {
			fmt.Printf("❌", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(haltCmd)
}
