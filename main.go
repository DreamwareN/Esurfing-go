package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/DreamwareN/Esurfing-go/client"
	"github.com/DreamwareN/Esurfing-go/config"
	"github.com/DreamwareN/Esurfing-go/utils"
)

var ClientList []*client.Client

func main() {
	var err error
	f1, err := os.Create("./prof/p1.prof")

	err = pprof.StartCPUProfile(f1)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	var configFilePath = flag.String("c", "config.json", "config file path")
	flag.Parse()

	fmt.Println("esurfing client v25.10.15")

	ClientList = []*client.Client{}

	log.Println("reading config")

	err = config.LoadConfig(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("load %d from:%s", len(config.List), *configFilePath)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for _, c := range config.List {
		wg.Add(1)
		err := RunNewClient(c, &wg, ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sigs

	log.Println("stoping all clients")
	cancel()
	for _, cl := range ClientList {
		go func() {
			cl.HttpClient.CloseIdleConnections()
		}()
	}
	wg.Wait()
}

func RunNewClient(c *config.Config, wg *sync.WaitGroup, ctx context.Context) error {
	transport, err := NewHttpTransport(c)
	if err != nil {
		return errors.New(fmt.Errorf("failed to create transport: %w", err).Error())
	}

	cl := &client.Client{
		Config: c,
		Ctx:    ctx,
		Wg:     wg,
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

	go cl.Start()
	return nil
}

func getLogPrefix(c *config.Config) string {
	return "User:" + c.Username + " BindDevice:" + func() string {
		if c.BindDevice != "" {
			return c.BindDevice
		}
		return "SystemDefault"
	}() + " "
}

func NewHttpTransport(c *config.Config) (http.RoundTripper, error) {
	if c.BindDevice != "" {
		ip, err := utils.GetInterfaceIP(c.BindDevice)
		if err != nil {
			return nil, errors.New(fmt.Errorf("failed to get interface IP: %w", err).Error())
		}

		localAddr := &net.TCPAddr{IP: net.ParseIP(ip)}
		return &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: localAddr,
				Resolver:  GetResolver(c),
			}).DialContext,
		}, nil
	} else {
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
