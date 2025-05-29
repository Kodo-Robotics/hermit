/*
Copyright © 2025 Kodo Robotics
*/
package boxcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Kodo-Robotics/hermit/pkg/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all added Hermit boxes",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := core.LoadRegistry()
		if err != nil {
			fmt.Println("❌ Failed to load box registry:", err)
			return
		}

		if len(reg) == 0 {
			fmt.Println("📦 No boxes have been added yet.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPATH\tADDED AT")

		for alias, info := range reg {
			fmt.Fprintf(w, "%s\t%s\t%s\n", alias, info.Path, info.AddedAt)
		}

		w.Flush()
	},
}

func init() {
	BoxCmd.AddCommand(listCmd)
}
