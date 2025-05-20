/*
Copyright © 2025 Kodo-Robotics

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
)

var vmName string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Virtualbox VM",
	Run: func(cmd *cobra.Command, args []string) {
		err := virtualbox.CreateVM(vmName)
		if err != nil {
			fmt.Println("❌", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&vmName, "name", "n", "hermit-vm", "Name of the VM")
}
