/*
Copyright © 2025 Kodo Robotics
*/
package core

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func RunSSH(user string, port int) error {
	fmt.Printf("🔐 Connecting to %s@localhost:%d ...\n", user, port)

	keyPath, err := GetOrInstallDefaultSSHKey()
	if err != nil {
		return fmt.Errorf("failed to prepare SSH key: %v", err)
	}

	cmd := exec.Command("ssh", "-p", strconv.Itoa(port),
		fmt.Sprintf("%s@localhost", user),
		"-i", keyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("SSH failed: %w", err)
	}

	fmt.Println("👋 SSH session ended.")
	return nil
}
