package client

import (
	"context"
	"errors"
	"github.com/DreamwareN/Esurfing-go/cipher"
	"github.com/DreamwareN/Esurfing-go/config"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	Conf       *config.Config
	Log        *log.Logger
	HttpClient *http.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
	WaitGroup  *sync.WaitGroup
	cipher     cipher.Cipher

	HeartbeatInterval  time.Duration
	IsConnectedAtFirst int
	IsCheckingNetwork  atomic.Int32
	IsRunning          atomic.Int32
	IsLogin            atomic.Int32

	UserIP     string
	AcIP       string
	Domain     string
	Area       string
	SchoolID   string
	ClientID   uuid.UUID
	Hostname   string
	MacAddress string
	Ticket     string

	AlgoID           string
	IndexURL         string
	TicketURL        string
	AuthURL          string
	KeepURL          string
	TermURL          string
	FirstRedirectURL string
}

var ErrMaxReTryReach = errors.New("max retry reached")

func (cl *Client) Run() {
	cl.Log.Println("Starting client...")
	defer cl.WaitGroup.Done()
	defer cl.Exit()

	networkCheckTicker := time.NewTicker(time.Millisecond * time.Duration(cl.Conf.NetworkCheckIntervalMS))
	defer networkCheckTicker.Stop()

	for {
		select {
		case <-networkCheckTicker.C:
			//todo: useless impl:cl.IsCheckingNetwork
			go func() {
				if cl.IsCheckingNetwork.Load() == 0 {
					cl.IsCheckingNetwork.Store(1)
					defer cl.IsCheckingNetwork.Store(0)
				} else if cl.IsCheckingNetwork.Load() == 1 {
					return
				}

				if err := cl.CheckNetworkStatus(); err != nil {
					if errors.Is(ErrMaxReTryReach, err) {
						cl.Log.Printf("%s exit", cl.Conf.AuthUsername)
						return
					}
					cl.Log.Printf("Network check failed: %v", err)
				}

			}()
		case <-cl.Ctx.Done():
			return
		}
	}
}

func (cl *Client) Exit() {
	if cl.IsRunning.Load() == 1 && cl.IsLogin.Load() == 1 {
		_ = cl.Logout()
	}
	cl.Log.Println("exit")
}
