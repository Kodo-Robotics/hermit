/*
Copyright © 2025 Kodo Robotics

*/
package virtualbox

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func CleanupStaleVBoxFile(vmName string) error {
	var basePath string

	if runtime.GOOS == "windows" {
		basePath = filepath.Join(os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH"), "VirtualBox VMs", vmName)
	} else {
		basePath = filepath.Join(os.Getenv("HOME"), "VirtualBox VMs", vmName)
	}

	vboxFile := filepath.Join(basePath, vmName + ".vbox")

	if _, err := os.Stat(vboxFile); err == nil {
		fmt.Println("⚠️ Removing stale .vbox file:", vboxFile)
		if err := os.Remove(vboxFile); err != nil {
			return fmt.Errorf("failed to remove stale .vbox: %v", err)
		}
	}

	return nil
}

func removeVBoxVMFolder(vmName string) error {
	var basePath string

	if runtime.GOOS == "windows" {
		basePath = filepath.Join(os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH"), "VirtualBox VMs", vmName)
	} else {
		basePath = filepath.Join(os.Getenv("HOME"), "VirtualBox VMs", vmName)
	}

	if _, err := os.Stat(basePath); err == nil {
		fmt.Println("⚠️ Removing pre-existing VM folder:", basePath)
		if err := os.RemoveAll(basePath); err != nil {
			return fmt.Errorf("failed to remove existing VM folder: %v", err)
		}
	}

	return nil
}