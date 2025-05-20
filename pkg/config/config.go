/*
Copyright © 2025 Kodo Robotics

*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type HermitConfig struct {
	Box				string	`json:"box"`
	Provider		string	`json:"provider"`
	Name			string	`json:"name"`
	CPUs			int		`json:"cpus"`
	Memory			int		`json:"memory"`
	DiskSizeMB		int		`json:"disk_size_mb"`
	ForwardedPorts	[]Port	`json:"forwarded_ports"`
}

type Port struct {
	Guest int `json:"guest"`
	Host  int `json:"host"`
}

func LoadConfig() (*HermitConfig, error) {
	data, err := os.ReadFile("hermit.json")
	if err != nil {
		return nil, fmt.Errorf("could not read hermit.json: %w", err)
	}

	var cfg HermitConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config format: %w", err)
	}

	return &cfg, nil
}