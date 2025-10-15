package config

import (
	"encoding/json"
	"errors"
	"math"
	"os"
)

type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckInterval int    `json:"check_interval"`
	RetryInterval int    `json:"retry_interval"`
	BindDevice    string `json:"bind_device"`
	DnsAddress    string `json:"dns_address"`
}

var List []*Config

func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("config file does not exist: " + configPath)
		}
		return err
	}
	err = json.Unmarshal(file, &List)
	if err != nil {
		return errors.New("load config file error: " + err.Error())
	}
	return check(List)
}

func check(cs []*Config) error {
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

		if c.Username == "" || c.Password == "" {
			return errors.New("username or password is empty")
		}
	}

	return nil
}
