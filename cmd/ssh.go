/*
Copyright © 2025 Kodo Robotics
*/
package cmd

import (
	"fmt"

	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/core"
	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into the Hermit VM (must have guest port 22 forwarded)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		user := "vagrant"
		port := 2222

		for _, p := range cfg.ForwardedPorts {
			if p.Guest == 22 {
				port = p.Host
				break
			}
		}

		if err := core.RunSSH(user, port); err != nil {
			fmt.Println("❌ SSH failed:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
