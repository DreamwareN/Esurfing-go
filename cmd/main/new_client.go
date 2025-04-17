package main

import (
	"context"
	"fmt"
	"github.com/DreamwareN/Esurfing-go/client"
	"github.com/DreamwareN/Esurfing-go/config"
	"log"
	"net/http"
	"os"
	"sync"
)

func RunNewClient(c *config.Config, wg *sync.WaitGroup) error {
	SetDefaultConfig(c)

	transport, err := NewHttpTransport(c)
	if err != nil {
		return fmt.Errorf("failed to create transport: %w", err)
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
			//这里有待商榷 默认的timeout太大了
			//Timeout:   time.Second * 3,
			Transport: transport,
		},
		AlgoID: "00000000-0000-0000-0000-000000000000",
	}

	cl.Log = log.New(
		os.Stdout,
		fmt.Sprintf("[User:%s Interface:%s] ", cl.Conf.AuthUsername, cl.Conf.NetworkInterfaceID),
		log.LstdFlags|log.Lmsgprefix,
	)

	ClientList = append(ClientList, cl)

	go cl.Run()
	return nil
}
