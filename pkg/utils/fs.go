/*
Copyright © 2025 Kodo Robotics
*/
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func FixKeyPermissions(path string) error {
	if runtime.GOOS == "windows" {
		user := os.Getenv("USERNAME")
		cmd1 := exec.Command("icacls", path, "/inheritance:r")
		cmd2 := exec.Command("icacls", path, "/grant:r", fmt.Sprintf("%s:R", user))

		if out, err := cmd1.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to remove inheritance: %v\n%s", err, string(out))
		}
		if out, err := cmd2.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set read-only for %s: %v\n%s", user, err, string(out))
		}
	} else {
		// Unix systems
		if err := os.Chmod(path, 0600); err != nil {
			return fmt.Errorf("failed to chmod 600: %v", err)
		}
	}

	return nil
}