/*
Copyright © 2025 Kodo Robotics
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Kodo-Robotics/hermit/pkg/config"
	"github.com/Kodo-Robotics/hermit/pkg/utils"
	"github.com/Kodo-Robotics/hermit/pkg/virtualbox"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Boot up the Hermit VM",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("❌ Error reading hermit.json:", err)
			return
		}

		// Check if VM already exists
		_, err = virtualbox.GetVMState(cfg.Name)
		if err == nil {
			if err := virtualbox.StartVM(cfg.Name); err != nil {
				fmt.Println("❌ Failed to start VM:", err)
			}
			return
		} else {
			// VM not found - check for stale vbox
			if err := virtualbox.CleanupStaleVBoxFile(cfg.Name); err != nil {
				fmt.Println("⚠️ Warning: could not clean stale .vbox file:", err)
			}
		}

		boxDir := utils.GetBoxPath(cfg.Box)
		ovfPath, err := utils.FindOVF(boxDir)
		if err != nil {
			fmt.Println("❌", err)
			fmt.Println("👉 Run `hermit box add <path>.box` to install the box.")
			return
		}

		fmt.Printf("📦 Importing VM '%s' from box: %s\n", cfg.Name, cfg.Box)
		if err := virtualbox.ImportOVF(ovfPath, cfg.Name); err != nil {
			fmt.Println("❌ Failed to import OVF:", err)
			return
		}

		fmt.Println("⚙️ Applying CPU and memory settings...")
		if err := virtualbox.ModifyVM(cfg.Name, cfg.Memory, cfg.CPUs); err != nil {
			fmt.Println("⚠️ Failed to modify VM settings:", err)
		}

		fmt.Println("🌐 Configuring networking...")
		selectNetworkAdapter(&cfg.Network)
		net := cfg.Network
		if err := virtualbox.ConfigureNetworking(cfg.Name, net.Mode, net.BridgeAdapter, net.HostOnlyAdapter); err != nil {
			fmt.Println("❌ Failed to configure networking:", err)
			return
		}

		for _, port := range net.ForwardedPorts {
			fmt.Printf("🔁 Forwarding host:%d -> guest:%d\n", port.Host, port.Guest)

			// Delete existing rule if exists
			_ = virtualbox.DeletePortForwardRule(cfg.Name, port.Guest)

			if err := virtualbox.AddPortForward(cfg.Name, port.Guest, port.Host); err != nil {
				fmt.Printf("⚠️ Failed to add port forward: %v\n", err)
			}
		}

		fmt.Println("🚀 Starting VM...")
		if err := virtualbox.StartVM(cfg.Name); err != nil {
			fmt.Println("❌ Error starting VM:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func selectNetworkAdapter(net *config.NetworkConfig) {
	reader := bufio.NewReader(os.Stdin)

	switch net.Mode {
	case "bridged":
		if net.BridgeAdapter == "" {
			adapters, _ := virtualbox.ListBridgeAdapters()
			fmt.Println("🌐 Select a bridged adapter:")
			for i, a := range adapters {
				fmt.Printf("  [%d] %s\n", i+1, a)
			}
			fmt.Print("Enter number: ")
			input, _ := reader.ReadString('\n')
			i := 0
			fmt.Sscanf(strings.TrimSpace(input), "%d", &i)
			if i > 0 && i <= len(adapters) {
				net.BridgeAdapter = adapters[i-1]
			}
		}

	case "hostonly":
		if net.HostOnlyAdapter == "" {
			adapters, _ := virtualbox.ListHostOnlyAdapters()
			fmt.Println("🔒 Select a host-only adapter:")
			for i, a := range adapters {
				fmt.Printf("  [%d] %s\n", i+1, a)
			}
			fmt.Print("Enter number: ")
			input, _ := reader.ReadString('\n')
			i := 0
			fmt.Sscanf(strings.TrimSpace(input), "%d", &i)
			if i > 0 && i <= len(adapters) {
				net.HostOnlyAdapter = adapters[i-1]
			}
		}
	}
}
