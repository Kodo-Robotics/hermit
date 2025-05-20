/*
Copyright © 2025 Kodo-Robotics

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Hermit environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌱 Hermit environment initialized")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
