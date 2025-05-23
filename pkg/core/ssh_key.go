/*
Copyright © 2025 Kodo Robotics
*/
package core

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/insecure_private_key
var vagrantPrivateKey string

func GetOrInstallDefaultSSHKey() (string, error) {
	keyPath := filepath.Join(".hermit", "insecure_private_key")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		if err := os.MkdirAll(".hermit", 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(keyPath, []byte(vagrantPrivateKey), 0600); err != nil {
			return "", err
		}
		fmt.Println("🔑 Installed default Vagrant SSH key in .hermit/")
	}

	return keyPath, nil
}
