package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckInterval int    `json:"check_interval"`
	RetryInterval int    `json:"retry_interval"`
	BindInterface string `json:"bind_interface"`
	DnsAddress    string `json:"dns_address"`
}

var Configs []*Config

func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("config file does not exist: " + configPath)
		}
		return err
	}
	err = json.Unmarshal(file, &Configs)
	if err != nil {
		return errors.New("load config file error: " + err.Error())
	}
	return nil
}
