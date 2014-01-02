package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Root       string `json:"root"`
	MountPoint string `json:"mount_point"`
}

func parseConfig(cfgfile string) (*config, error) {
	var Cfg config
	jsonStr := make([]byte, 2048)

	file, err := os.Open(cfgfile)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	n, err := file.Read(jsonStr);
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonStr[:n], &Cfg); err != nil {
		return nil, err
	}
	return &Cfg, nil
}
