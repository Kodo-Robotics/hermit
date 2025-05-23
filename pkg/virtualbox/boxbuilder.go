/*
Copyright © 2025 Kodo Robotics
*/
package virtualbox

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExportVM(vmName string, outputDir string) (string, error) {
	ovfPath := filepath.Join(outputDir, "box.ovf")
	err := runVBoxManage("export", vmName, "--output", ovfPath)
	if err != nil {
		return "", fmt.Errorf("failed to export VM: %v", err)
	}
	return ovfPath, nil
}

func FindDiskFile(outputDir string) (string, error) {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".vmdk") {
			return filepath.Join(outputDir, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("no .vmdk file found in %s", outputDir)
}

func CreateBoxArchive(outputPath string, files map[string]string) error {
	boxFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer boxFile.Close()

	tw := tar.NewWriter(boxFile)
	defer tw.Close()

	for tarName, filePath := range files {
		if err := addFileToTar(tw, filePath, tarName); err != nil {
			return err
		}
	}
	return nil
}

func addFileToTar(tw *tar.Writer, filePath string, tarName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: tarName,
		Mode: 0644,
		Size: stat.Size(),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, f)
	return err
}
