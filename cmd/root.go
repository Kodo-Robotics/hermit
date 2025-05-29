/*
Copyright © 2025 Kodo Robotics
*/
package cmd

import (
	"os"

	"github.com/Kodo-Robotics/hermit/cmd/boxcmd"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hermit",
	Short: "Fast, minimal, and cross-platform CLI tool to manage virtual development environments",
	Long: `Hermit is a blazing-fast, minimal, and extensible CLI tool to manage virtual development environments.
Built in Go as a modern alternative to Vagrant, Hermit focuses on performance, simplicity, and cross-platform support.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(boxcmd.BoxCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
