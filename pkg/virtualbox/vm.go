/*
Copyright © 2025 Kodo Robotics

*/

package virtualbox

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ImportOVF(ovfPath string, desiredName string) error {
	originalName, err := ExtractVMNameFromOVF(ovfPath)
	if err != nil {
		return fmt.Errorf("failed to parse OVF: %v", err)
	}

	if err := runVBoxManage("import", ovfPath); err != nil {
		return fmt.Errorf("import failed: %v", err)
	}

	if originalName != desiredName {
		_ = removeVBoxVMFolder(desiredName)
		if err := runVBoxManage("modifyvm", originalName, "--name", desiredName); err != nil {
			return fmt.Errorf("failed to rename VM: %v", err)
		}
	}
	return nil
}

func GetVMState(vmName string) (string, error) {
	cmd := exec.Command("VBoxManage", "showvminfo", vmName, "--machinereadable")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("VM not found or error checking state: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VMState=") {
			state := strings.Trim(strings.Split(line, "=")[1], "\"")
			return state, nil
		}
	}

	return "", fmt.Errorf("VMState not found")
}

func ModifyVM(vmName string, memory int, cpus int) error {
	args := []string{"modifyvm", vmName}
	if memory > 0 {
		args = append(args, "--memory", fmt.Sprintf("%d", memory))
	}
	if cpus > 0 {
		args = append(args, "--cpus", fmt.Sprintf("%d", cpus))
	}

	return runVBoxManage(args...)
}

func AddPortForward(vmName string, guestPort, hostPort int) error {
	rule := fmt.Sprintf("fwd%d,tcp,,%d,,%d", guestPort, hostPort, guestPort)
	return runVBoxManage("modifyvm", vmName, "--natpf1", rule)
}

func StartVM(vmName string) error {
	state, err := GetVMState(vmName)
	if err != nil {
		return fmt.Errorf("VM '%s' not found: %v", vmName, err)
	}

	switch state {
	case "running":
		fmt.Println("✅ VM is already running.")
		return nil

	case "poweroff", "saved", "aborted":
		fmt.Println("🔁 VM exists. Starting...")
		err := runVBoxManage("startvm", vmName, "--type", "headless")
		if err != nil {
			return fmt.Errorf("failed to start VM: %v", err)
		}

		// Wait for VM to be running
		fmt.Print("⏳ Waiting for VM to start")
		timeout := time.After(30 * time.Second)
		tick := time.Tick(1 * time.Second)

		for {
			select {
			case <-timeout:
				return fmt.Errorf("\n⏰ Timeout waiting for VM to start")
			case <-tick:
				current, _ := GetVMState(vmName)
				if current == "running" {
					fmt.Println("\n✅ VM is now running.")
					return nil
				}
				fmt.Print(".")
			}
		}

	default:
		return fmt.Errorf("🛑 VM is in unsupported state: %s", state)
	}
}

func HaltVM(vmName string) error {
	state, err := GetVMState(vmName)
	if err != nil {
		return fmt.Errorf("could not determine VM state: %v", err)
	}
	if state == "poweroff" {
		fmt.Println("⏹️ VM is already stopped.")
		return nil
	}

	// Send graceful shutdown
	fmt.Println("🛑 Sending ACPI shutdown signal...")
	err = runVBoxManage("controlvm", vmName, "acpipowerbutton")
	if err != nil {
		fmt.Println("⚠️ ACPI shutdown failed. Trying hard poweroff...")
		err = runVBoxManage("controlvm", vmName, "poweroff")
		if err != nil {
			return fmt.Errorf("failed to force shutdown: %v", err)
		}
	}

	// Wait for VM to shutdown
	fmt.Print("⏳ Waiting for VM to shut down")
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("\n⏰ Timeout waiting for VM to shut down")
		case <-tick:
			current, _ := GetVMState(vmName)
			if current == "poweroff" || current == "aborted" {
				fmt.Println("\n✅ VM is powered off.")
				return nil
			}
			fmt.Print(".")
		}
	}
}

func DestroyVM(vmName string, deleteDisks bool) error {
	state, err := GetVMState(vmName)
	if err != nil {
		return fmt.Errorf("VM '%s' not found or already removed", vmName)
	}

	if state == "running" {
		fmt.Println("🛑 VM is running, stopping before destroy...")
		if err := HaltVM(vmName); err != nil {
			fmt.Printf("⚠️ Failed to stop VM before destroy: %v\n", err)
		}
	}

	args := []string{"unregistervm", vmName}
	if deleteDisks {
		args = append(args, "--delete")
	}
	return runVBoxManage(args...)
}

func ConfigureNetworking(vmName string, netMode string, bridgeAdapter string, hostOnlyAdapter string) error {
	if err := runVBoxManage("modifyvm", vmName, "--nic1", "nat"); err != nil {
		return fmt.Errorf("failed to configure NIC1 as NAT: %v", err)
	}

	switch netMode {
	case "", "none":
		return nil

	case "bridged":
		if bridgeAdapter == "" {
			return fmt.Errorf("bridge_adapter must be set for bridged mode")
		}
		return runVBoxManage("modifyvm", vmName, "--nic2", "bridged", "--bridgeadapter2", bridgeAdapter)

	case "hostonly":
		if hostOnlyAdapter == "" {
			return fmt.Errorf("hostonly_adapter must be set for hostonly mode")
		}
		return runVBoxManage("modifyvm", vmName, "--nic2", "hostonly", "--hostonlyadapter2", hostOnlyAdapter)

	default:
		return fmt.Errorf("invalid network mode: %s", netMode)
	}
}

func runVBoxManage(args ...string) error {
	cmd := exec.Command("VBoxManage", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("VBoxManage error: %v\n%s", err, string(out))
	}

	return nil
}
