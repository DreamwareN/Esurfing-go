package main

import (
	"fmt"
	"github.com/DreamwareN/Esurfing-go/client/utils"
	"github.com/DreamwareN/Esurfing-go/config"
	"math"
	"net"
	"net/http"
)

func SetDefaultConfig(c *config.Config) {
	if c.NetworkCheckIntervalMS == 0 {
		c.NetworkCheckIntervalMS = 3000
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 5
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = math.MaxInt64
	}
	if c.RetryDelayMS == 0 {
		c.RetryDelayMS = 1000
	}
}

func NewHttpTransport(c *config.Config) (http.RoundTripper, error) {
	if c.NetworkInterfaceID == "" {
		return http.DefaultTransport, nil
	}

	ip, err := utils.GetInterfaceIP(c.NetworkInterfaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface IP: %w", err)
	}

	localAddr := &net.TCPAddr{IP: net.ParseIP(ip)}
	return &http.Transport{
		DialContext: (&net.Dialer{
			LocalAddr: localAddr,
		}).DialContext,
	}, nil
}
