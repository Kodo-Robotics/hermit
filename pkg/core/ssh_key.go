/*
Copyright © 2025 Kodo Robotics
*/
package core

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kodo-Robotics/hermit/pkg/utils"
)

//go:embed assets/insecure_private_key
var vagrantPrivateKey string

func GetOrInstallDefaultSSHKey() (string, error) {
	keyPath := filepath.Join(utils.GetHermitRoot(), "insecure_private_key")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		if err := os.MkdirAll(utils.GetHermitRoot(), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(keyPath, []byte(vagrantPrivateKey), 0600); err != nil {
			return "", err
		}

		if err := utils.FixKeyPermissions(keyPath); err != nil {
			fmt.Println("⚠️ Warning: could not restrict SSH key permissions:", err)
		}

		fmt.Println("🔑 Installed default Vagrant SSH key in .hermit/")
	}

	return keyPath, nil
}
