/*
Copyright © 2025 Kodo Robotics

*/

package virtualbox

import (
	"fmt"
	"os/exec"
	"strings"
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
	return runVBoxManage("startvm", vmName, "--type", "headless")
}

func HaltVM(vmName string) error {
	return runVBoxManage("controlvm", vmName, "acpipowerbutton")
}

func runVBoxManage(args ...string) error {
	cmd := exec.Command("VBoxManage", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("VBoxManage error: %v\n%s", err, string(out))
	}

	return nil
}