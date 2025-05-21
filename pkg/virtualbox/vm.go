/*
Copyright © 2025 Kodo Robotics

*/

package virtualbox

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func runVBoxManage(args ...string) error {
	cmd := exec.Command("VBoxManage", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("VBoxManage error: %v\n%s", err, string(out))
	}

	return nil
}

func CreateAndStartVM(name string, memory int, cpus int, vram int, graphicsController string, vdiPath string) error {
	diskPath, err := filepath.Abs(vdiPath)
	if err != nil {
		return fmt.Errorf("failed to resolve VDI path: %v", err)
	}

	fmt.Println("📦 Creating VM...")
	if err := runVBoxManage("createvm", "--name", name, "--register"); err != nil {
		return err
	}

	fmt.Println("🔌 Adding SATA controller...")
	if err := runVBoxManage("storagectl", name, "--name", "SATA Controller", "--add", "sata", "--controller", "IntelAhci"); err != nil {
		return err
	}

	fmt.Println("📎 Attaching virtual hard disk...")
	if err := runVBoxManage("storageattach", name, "--storagectl", "SATA Controller", "--port", "0", "--device", "0", "--type", "hdd", "--medium", diskPath); err != nil {
		return err
	}

	fmt.Println("⚙️ Configuring VM resources...")
	if err := runVBoxManage("modifyvm", name, "--memory", fmt.Sprintf("%d", memory), "--cpus", fmt.Sprintf("%d", cpus)); err != nil {
		return err
	}

	if err := runVBoxManage("modifyvm", name, "--vram", fmt.Sprintf("%d", vram), "--graphicscontroller", graphicsController); err != nil {
		return err
	}

	fmt.Println("🎬 Starting VM in headless mode...")
	if err := runVBoxManage("startvm", name, "--type", "headless"); err != nil {
		return err
	}

	fmt.Println("✅ VM started successfully.")
	return nil
}