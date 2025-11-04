package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var clients []*Client
var wg sync.WaitGroup

func main() {
	var err error
	var configFilePath = flag.String("c", "config.json", "config file path")
	flag.Parse()

	log.Println("esurfing client v25.11.4")
	log.Println("reading config")

	err = LoadConfig(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("load %d from:%s", len(Configs), *configFilePath)

	for _, c := range Configs {
		client, err := NewClient(c)
		if err != nil {
			log.Fatal(err)
		}

		clients = append(clients, client)

		go client.Start()

		wg.Add(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-signalChannel

	log.Println("stoping all clients")

	for _, client := range clients {
		client.Cancel()
	}

	wg.Wait()
	log.Println("exit")
}
