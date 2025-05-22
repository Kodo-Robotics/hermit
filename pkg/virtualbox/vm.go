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

func ImportOVF(ovfPath string, vmName string) error {
	return runVBoxManage("import", ovfPath, "--vsys", "0", "--vmname", vmName)
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

func runVBoxManage(args ...string) error {
	cmd := exec.Command("VBoxManage", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("VBoxManage error: %v\n%s", err, string(out))
	}

	return nil
}