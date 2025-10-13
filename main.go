package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DreamwareN/Esurfing-go/client"
	"github.com/DreamwareN/Esurfing-go/config"
	"github.com/DreamwareN/Esurfing-go/errs"
	"github.com/DreamwareN/Esurfing-go/utils"
)

var ClientList []*client.Client

func main() {
	var configFilePath = flag.String("c", "config.json", "config file path")
	flag.Parse()

	fmt.Println("ESurfing-go ver:", Version)

	ClientList = make([]*client.Client, 0)
	wg := sync.WaitGroup{}

	log.Println("reading config...")

	err := config.LoadConfig(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d account(s) loaded from config: %s", len(config.ConfigList), *configFilePath)

	for _, c := range config.ConfigList {
		wg.Add(1)
		err := RunNewClient(c, &wg)
		if err != nil {
			log.Fatal(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sigs

	log.Println("shutting down...")
	for _, cl := range ClientList {
		go func() {
			cl.Cancel()
			cl.HttpClient.CloseIdleConnections()
		}()
	}
	wg.Wait()
	log.Println("bye")
}

func RunNewClient(c *config.Config, wg *sync.WaitGroup) error {
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
	if c.DnsAddress == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, "udp", c.DnsAddress)
		},
	}
}
