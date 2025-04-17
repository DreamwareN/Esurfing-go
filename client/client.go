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

	isRunning atomic.Int32
	isLogin   atomic.Int32

	UserIP     string
	AcIP       string
	Domain     string
	Area       string
	SchoolID   string
	ClientID   uuid.UUID
	Hostname   string
	MacAddress string
	Ticket     string
	Key        string

	AlgoID    string
	IndexURL  string
	TicketURL string
	AuthURL   string
	KeepURL   string
	TermURL   string

	FirstRedirectURL  string
	SecondRedirectURL string

	HeartbeatInterval time.Duration

	//todo: remove this later
	//《网络是否已连接》
	isConnectedAtFirst int
}

var ErrMaxReTryReach = errors.New("max retry reached")

func (cl *Client) Run() {
	cl.Log.Println("Starting client...")
	defer cl.WaitGroup.Done()
	defer cl.exit()

	networkCheckTicker := time.NewTicker(time.Millisecond * time.Duration(cl.Conf.NetworkCheckIntervalMS))
	defer networkCheckTicker.Stop()

	for {
		select {
		case <-networkCheckTicker.C:
			if err := cl.checkNetworkStatus(); err != nil {
				if errors.Is(ErrMaxReTryReach, err) {
					cl.Log.Printf("%s exit", cl.Conf.AuthUsername)
					return
				}
				cl.Log.Printf("Network check failed: %v", err)
			}

		case <-cl.Ctx.Done():
			return
		}
	}
}

func (cl *Client) exit() {
	if cl.isRunning.Load() == 1 && cl.isLogin.Load() == 1 {
		if err := cl.logout(); err != nil {
			//这里永远都是timeout 不用处理后续
			//cl.Log.Println("logout failed: ", err.Error())
		}
	}
	cl.Log.Println("exit")
}
