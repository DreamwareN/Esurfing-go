package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	AuthUsername           string `json:"username"`
	AuthPassword           string `json:"password"`
	NetworkCheckIntervalMS int    `json:"network_check_interval_ms"`
	MaxRetries             int    `json:"max_retry"`
	RetryDelayMS           int    `json:"retry_delay_ms"`
	NetworkInterfaceID     string `json:"network_interface_id"`
}

var Conf []*Config

func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("config file does not exist: " + configPath)
		}
		return err
	}
	err = json.Unmarshal(file, &Conf)
	if err != nil {
		return errors.New("failed to load config file: " + err.Error())
	}
	return nil
}
