/*
Copyright © 2025 Kodo Robotics
*/
package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Kodo-Robotics/hermit/pkg/utils"
)

type BoxInfo struct {
	Path    string `json:"path"`
	AddedAt string `json:"added_at"`
}

type Registry map[string]BoxInfo

func LoadRegistry() (Registry, error) {
	path := filepath.Join(utils.GetHermitBoxPath(), "box_registry.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(Registry), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func SaveRegistry(reg Registry) error {
	path := filepath.Join(utils.GetHermitBoxPath(), "box_registry.json")
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddBox(name, path string) error {
	reg, err := LoadRegistry()
	if err != nil {
		return err
	}

	reg[name] = BoxInfo{
		Path:    path,
		AddedAt: time.Now().UTC().Format(time.RFC3339),
	}

	return SaveRegistry(reg)
}

func RemoveBox(name string) error {
	reg, err := LoadRegistry()
	if err != nil {
		return err
	}

	if _, ok := reg[name]; !ok {
		return fmt.Errorf("box '%s' not found in registry", name)
	}

	delete(reg, name)
	return SaveRegistry(reg)
}

func GetBox(name string) (BoxInfo, error) {
	reg, err := LoadRegistry()
	if err != nil {
		return BoxInfo{}, err
	}

	info, ok := reg[name]
	if !ok {
		return BoxInfo{}, fmt.Errorf("box '%s' not found", name)
	}

	return info, nil
}
