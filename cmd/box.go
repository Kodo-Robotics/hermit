/*
Copyright © 2025 Kodo Robotics

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "Manage Hermit boxes",
}

func init() {
	rootCmd.AddCommand(boxCmd)
}
