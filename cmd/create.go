/*
Copyright © 2025 Kodo-Robotics

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var vmName string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Virtualbox VM",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("📦 Creating VirtualBox VM: %s\n", vmName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&vmName, "name", "n", "hermit-vm", "Name of the VM")
}
