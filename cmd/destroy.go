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

var deleteDisks bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Unregister and optionally delete the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()

		if !deleteDisks {
			fmt.Println("💡 Use '--delete' next time to fully remove disk and .vbox file.")
		}

		fmt.Printf("🔥 Destroying VM '%s'...\n", cfg.Name)
		err = virtualbox.DestroyVM(cfg.Name, deleteDisks)
		if err != nil {
			fmt.Println("❌ Failed to destroy VM:", err)
			return
		}
		fmt.Println("✅ VM destroyed.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().BoolVarP(&deleteDisks, "delete", "d", false, "Delete all virtual disk files")
}
