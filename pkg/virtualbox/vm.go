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

func CreateAndStartVM(name string, memory int, cpus int, diskSize int, isoPath string) error {
	diskPath := filepath.Join(".", fmt.Sprintf("%s.vdi", name))

	fmt.Println("📦 Creating VM...")
	if err := runVBoxManage("createvm", "--name", name, "--register"); err != nil {
		return err
	}

	fmt.Println("💾 Creating virtual hard disk...")
	if err := runVBoxManage("createhd", "--filename", diskPath, "--size", fmt.Sprintf("%d", diskSize)); err != nil {
		return err
	}

	fmt.Println("🔌 Adding SATA controller...")
	if err := runVBoxManage("storagectl", name, "--name", "SATA Controller", "--add", "sata", "--controller", "IntelAhci"); err != nil {
		return err
	}

	fmt.Println("📎 Attaching hard disk...")
	if err := runVBoxManage("storageattach", name, "--storagectl", "SATA Controller", "--port", "0", "--device", "0", "--type", "hdd", "--medium", diskPath); err != nil {
		return err
	}

	fmt.Println("📀 Attaching ISO...")
	if err := runVBoxManage("storagectl", name, "--name", "IDE Controller", "--add", "ide"); err != nil {
		return err
	}

	if err := runVBoxManage("storageattach", name, "--storagectl", "IDE Controller", "--port", "0", "--device", "0", "--type", "dvddrive", "--medium", isoPath); err != nil {
		return err
	}

	fmt.Println("⚙️ Configuring VM resources...")
	if err := runVBoxManage("modifyvm", name, "--memory", fmt.Sprintf("%d", memory), "--cpus", fmt.Sprintf("%d", cpus)); err != nil {
		return err
	}

	fmt.Println("🚀 Setting boot order...")
	if err := runVBoxManage("modifyvm", name, "--boot1", "dvd", "--boot2", "disk", "--boot3", "none"); err != nil {
		return err
	}

	fmt.Println("🎬 Starting VM in headless mode...")
	if err := runVBoxManage("startvm", name, "--type", "gui"); err != nil {
		return err
	}

	fmt.Println("✅ VM started successfully.")
	return nil
}