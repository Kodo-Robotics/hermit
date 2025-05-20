package virtualbox

import (
	"fmt"
	"os/exec"
	"github.com/google/uuid"
)

// CreateVM creates and registers a new VirtualBox VM with a unique name
func CreateVM(name string) error {
	if name == "" {
		name = fmt.Sprintf("hermit-%s", uuid.New().String()[:8])
	}

	// Create the VM
	createCmd := exec.Command("VBoxManage", "createvm", "--name", name, "--register")
	output, err := createCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create VM: %v\n%s", err, string(output))
	}

	fmt.Println("✅ VM created and registered:", name)

	modifyCmd := exec.Command("VBoxManage", "modifyvm", name,
		"--memory", "1024", "--cpus", "1", "--ostype", "Ubuntu_64")
	output, err = modifyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to modify VM: %v\n%s", err, string(output))
	}

	fmt.Println("⚙️ VM configured with 1 CPU, 1024MB RAM, Ubuntu 64-bit")

	return nil
}