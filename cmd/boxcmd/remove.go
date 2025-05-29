/*
Copyright © 2025 Kodo Robotics
*/
package boxcmd

import (
	"fmt"
	"os"

	"github.com/Kodo-Robotics/hermit/pkg/core"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a Hermit box",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		boxName := args[0]

		info, err := core.GetBox(boxName)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		if err := os.RemoveAll(info.Path); err != nil {
			fmt.Println("❌ Failed to delete box folder:", err)
			return
		}

		if err := core.RemoveBox(boxName); err != nil {
			fmt.Println("❌ Failed to update registry:", err)
			return
		}

		fmt.Printf("🗑️  Box '%s' removed successfully.\n", boxName)
	},
}

func init() {
	BoxCmd.AddCommand(removeCmd)
}
