/*
Copyright © 2025 Kodo Robotics
*/
package utils

import (
	"os"
	"path/filepath"
)

func GetHermitRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not detect user home directory")
	}
	return filepath.Join(home, ".hermit")
}

func GetBoxPath(boxName string) string {
	safe := filepath.Join(GetHermitRoot(), "boxes", sanitizeName(boxName))
	return safe
}

func sanitizeName(box string) string {
	return filepath.Clean(box)
}