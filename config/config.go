package config

import (
	"encoding/json"
	"errors"
	"math"
	"os"
)

type Config struct {
	AuthUsername           string `json:"1username"`
	AuthPassword           string `json:"1password"`
	NetworkCheckIntervalMS int    `json:"network_check_interval_ms"`
	MaxRetries             int    `json:"max_retry"`
	RetryDelayMS           int    `json:"retry_delay_ms"`
	OutBoundType           string `json:"out_bound_type"`
	NetworkInterfaceID     string `json:"network_interface_id"`
	NetworkBindAddress     string `json:"network_bind_address"`
	UseCustomDns           bool   `json:"use_custom_dns"`

	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckInterval int    `json:"check_interval"`
	RetryInterval int    `json:"retry_interval"`
	BindDevice    string `json:"bind_device"`
	DnsAddress    string `json:"dns_address"`
}

var ConfigList []*Config

func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("config file does not exist: " + configPath)
		}
		return err
	}
	err = json.Unmarshal(file, &ConfigList)
	if err != nil {
		return errors.New("failed to load config file: " + err.Error())
	}
	check(ConfigList)
	return nil
}

func check(cs []*Config) {
	for _, c := range cs {
		if c.CheckInterval <= 0 {
			c.CheckInterval = 3000
		}

		if c.RetryInterval == 0 {
			c.RetryInterval = 10000
		}

		if c.RetryInterval < 0 {
			c.RetryInterval = math.MaxInt32
		}
	}
}
