package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/DreamwareN/Esurfing-go/cipher"
	"github.com/DreamwareN/Esurfing-go/config"
	"github.com/DreamwareN/Esurfing-go/errs"
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

	HeartbeatInterval time.Duration
	IsCheckingNetwork atomic.Bool
	IsRunning         atomic.Bool
	IsLogin           atomic.Bool

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

	cl.CheckFirst()

	for {
		select {
		case <-networkCheckTicker.C:
			go func() {
				if cl.IsCheckingNetwork.Load() {
					return
				}

				cl.IsCheckingNetwork.Store(true)
				defer cl.IsCheckingNetwork.Store(false)

				if err := cl.CheckNetworkStatus(); err != nil {
					if errors.Is(ErrMaxReTryReach, err) {
						cl.Log.Println("Max auth retry reached")
						networkCheckTicker.Stop()
						return
					}
					cl.Log.Printf("Network check failed:%v", err)
				}
			}()
		case <-cl.Ctx.Done():
			return
		}
	}

}

func (cl *Client) Exit() {
	if cl.IsRunning.Load() && cl.IsLogin.Load() {
		_ = cl.Logout()
	}
	cl.Log.Println("Exit")
}

func (cl *Client) CheckFirst() {
	request, err := cl.GenerateGetRequest("http://connect.rom.miui.com/generate_204")
	if err != nil {
		return
	}

	resp, err := cl.HttpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		cl.Log.Println("The network has been connected")
	}
}

func (cl *Client) CheckNetworkStatus() error {
	request, err := cl.GenerateGetRequest("http://connect.rom.miui.com/generate_204")
	if err != nil {
		return errs.New(err.Error())
	}

	resp, err := cl.HttpClient.Do(request)
	if err != nil {
		return errs.New(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		cl.IsRunning.Store(true)
		return nil

	case http.StatusFound:
		cl.Log.Println("Authorization required")
		return cl.HandleRedirect(resp)

	default:
		return errs.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}

func (cl *Client) HandleRedirect(resp *http.Response) error {
	redirectURL := resp.Header.Get("Location")
	if redirectURL == "" {
		return errs.New("missing Location header in redirect")
	}

	retryCount := 0
	for retryCount < cl.Conf.MaxRetries || (cl.Conf.MaxRetries == 0 && retryCount == 0) {
		select {
		case <-cl.Ctx.Done():
			return errs.New("context canceled")
		default:
		}

		if err := cl.Authorization(redirectURL); err != nil {
			retryCount++
			cl.Log.Printf("Authorization attempt %d failed: %v", retryCount, err)
			select {
			case <-time.After(time.Duration(cl.Conf.RetryDelayMS) * time.Millisecond):
			case <-cl.Ctx.Done():
				return errs.New("context canceled")
			}
			continue
		}

		cl.Log.Println("Authorization succeeded")
		go cl.MaintainSession()
		return nil
	}

	cl.Log.Println("Reached max retry count: ", cl.Conf.MaxRetries)
	return ErrMaxReTryReach
}
