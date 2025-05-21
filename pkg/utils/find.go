package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindDiskImage(boxDir string) (string, error) {
	files, err := os.ReadDir(boxDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.Type().IsRegular() {
			continue
		}
		lowerName := strings.ToLower(file.Name())

		// AppleDouble macOS files
		if strings.HasPrefix(lowerName, "._") {
			continue
		}

		if strings.HasSuffix(lowerName, ".vdi") || strings.HasSuffix(lowerName, ".vmdk") {
			return filepath.Join(boxDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no .vdi or .vmdk file found in %s", boxDir)
}