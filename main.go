package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DreamwareN/Esurfing-go/client"
	"github.com/DreamwareN/Esurfing-go/config"
)

var ClientList []*client.Client

func main() {
	var configFilePath = flag.String("c", "config.json", "config file path")
	flag.Parse()

	fmt.Println("ESurfing-go ver:", Version)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ClientList = make([]*client.Client, 0)
	wg := sync.WaitGroup{}

	log.Println("reading config...")

	err := config.LoadConfig(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d account(s) loaded from config: %s", len(config.Conf), *configFilePath)

	for _, c := range config.Conf {
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
