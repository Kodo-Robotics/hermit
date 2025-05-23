/*
Copyright © 2025 Kodo Robotics
*/
package virtualbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func ListBridgeAdapters() ([]string, error) {
	cmd := exec.Command("VBoxManage", "list", "bridgedifs")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list bridged adapters: %v", err)
	}
	return parseAdapterNames(out), nil
}

func ListHostOnlyAdapters() ([]string, error) {
	cmd := exec.Command("VBoxManage", "list", "hostonlyifs")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list host-only adapters: %v", err)
	}
	return parseAdapterNames(out), nil
}

func CleanupStaleVBoxFile(vmName string) error {
	var basePath string

	if runtime.GOOS == "windows" {
		basePath = filepath.Join(os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH"), "VirtualBox VMs", vmName)
	} else {
		basePath = filepath.Join(os.Getenv("HOME"), "VirtualBox VMs", vmName)
	}

	vboxFile := filepath.Join(basePath, vmName+".vbox")

	if _, err := os.Stat(vboxFile); err == nil {
		fmt.Println("⚠️ Removing stale .vbox file:", vboxFile)
		if err := os.Remove(vboxFile); err != nil {
			return fmt.Errorf("failed to remove stale .vbox: %v", err)
		}
	}

	return nil
}

func parseAdapterNames(data []byte) []string {
	var names []string
	re := regexp.MustCompile(`(?m)^Name:\s+(.*)$`)
	matches := re.FindAllSubmatch(data, -1)
	for _, m := range matches {
		names = append(names, strings.TrimSpace(string(m[1])))
	}
	return names
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
