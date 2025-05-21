package main

import (
	"context"
	"fmt"
	"github.com/DreamwareN/Esurfing-go/client"
	"github.com/DreamwareN/Esurfing-go/config"
	"github.com/DreamwareN/Esurfing-go/errs"
	"github.com/DreamwareN/Esurfing-go/utils"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func RunNewClient(c *config.Config, wg *sync.WaitGroup) error {
	SetDefaultConfig(c)

	transport, err := NewHttpTransport(c)
	if err != nil {
		return errs.New(fmt.Errorf("failed to create transport: %w", err).Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	cl := &client.Client{
		Conf:      c,
		Ctx:       ctx,
		Cancel:    cancel,
		WaitGroup: wg,
		HttpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		},
		AlgoID: "00000000-0000-0000-0000-000000000000",
		Log: log.New(
			os.Stdout,
			getLogPrefix(c),
			log.LstdFlags|log.Lmsgprefix,
		),
	}

	ClientList = append(ClientList, cl)

	go cl.Run()
	return nil
}

func getLogPrefix(c *config.Config) string {
	if c.OutBoundType == "id" {
		return "[User:" + c.AuthUsername + " Interface:" + c.NetworkInterfaceID + "] "
	}

	if c.OutBoundType == "ip" {
		return "[User:" + c.AuthUsername + " BindAddr:" + c.NetworkBindAddress + "] "
	}

	return "[User:" + c.AuthUsername + "] "
}

func SetDefaultConfig(c *config.Config) {
	if c.NetworkCheckIntervalMS == 0 {
		c.NetworkCheckIntervalMS = 1000
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = math.MaxInt32
	}
	if c.RetryDelayMS == 0 {
		c.RetryDelayMS = 1000
	}
}

func NewHttpTransport(c *config.Config) (http.RoundTripper, error) {
	switch c.OutBoundType {
	case "id":
		if c.NetworkInterfaceID == "" {
			return nil, errs.New("empty network interface ID")
		}

		ip, err := utils.GetInterfaceIP(c.NetworkInterfaceID)
		if err != nil {
			return nil, errs.New(fmt.Errorf("failed to get interface IP: %w", err).Error())
		}

		localAddr := &net.TCPAddr{IP: net.ParseIP(ip)}
		return &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: localAddr,
				Resolver:  GetResolver(c),
			}).DialContext,
		}, nil

	case "ip":
		if c.NetworkBindAddress == "" {
			return nil, errs.New("empty network bind address")
		}

		ip := net.ParseIP(c.NetworkBindAddress)
		if ip == nil {
			return nil, errs.New("invalid network bind address")
		}

		localAddr := &net.TCPAddr{IP: ip}
		return &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: localAddr,
				Resolver:  GetResolver(c),
			}).DialContext,
		}, nil
	default:
		return &http.Transport{
			DialContext: (&net.Dialer{
				Resolver: GetResolver(c),
			}).DialContext,
		}, nil
	}
}

func GetResolver(c *config.Config) *net.Resolver {
	if !c.UseCustomDns {
		return net.DefaultResolver
	}

	if c.DnsAddress == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, "udp", c.DnsAddress+":53")
		},
	}
}
