/*
Copyright © 2025 Kodo Robotics
*/
package boxcmd

import (
	"github.com/spf13/cobra"
)

var BoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Manage Hermit boxes",
}

func init() {}
