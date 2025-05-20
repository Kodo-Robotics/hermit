// /*
// Copyright © 2025 Kodo-Robotics

// */
package cmd

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	"github.com/spf13/cobra"
// 	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
// 	"github.com/Kodo-Robotics/hermit/pkg/config"
// )

// var upCmd = &cobra.Command{
// 	Use:   "up",
// 	Short: "Boot up the Hermit VM",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		cfg, err := config.LoadConfig()
// 		if err != nil {
// 			fmt.Println("❌ Error reading hermit.json:", err)
// 			return
// 		}

// 		err = virtualbox.CreateAndStartVM(cfg.Name, cfg.Memory, cfg.CPUs, cfb.DiskSizeMB, cfg.ISOPath)
// 		if err != nil {
// 			fmt.Println("❌ Failed to bring up VM:", err)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(upCmd)
// }
