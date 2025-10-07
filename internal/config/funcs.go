package config

import (
	"encoding/json"
	"os"
)

func SetUser(cfg Config) error {
	if err := write(cfg); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	cfg := Config{}
	pFile, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(pFile)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	var data []byte
	pFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err = json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(pFile, data, 0o644); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	path += "/" + cfgFileName
	return path, err
}
