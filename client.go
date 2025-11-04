package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Config          *Config
	Log             *log.Logger
	HttpClient      *http.Client
	Ctx             context.Context
	Cancel          context.CancelFunc
	cipher          Cipher
	heartBeatTicker *time.Ticker

	UserIP     string
	AcIP       string
	Domain     string
	Area       string
	SchoolID   string
	ClientID   uuid.UUID
	Hostname   string
	MacAddress string
	Ticket     string
	AlgoID     string

	IndexUrl    string
	TicketUrl   string
	AuthUrl     string
	KeepUrl     string
	TermUrl     string
	RedirectUrl string
}

func NewClient(config *Config) (*Client, error) {
	if config.CheckInterval <= 0 {
		config.CheckInterval = 10000
	}

	if config.RetryInterval == 0 {
		config.RetryInterval = 10000
	}

	if config.RetryInterval < 0 {
		config.RetryInterval = math.MaxInt32
	}

	if config.Username == "" || config.Password == "" {
		return nil, errors.New("username or password is empty")
	}

	if config.BindInterface == "" {
		config.BindInterface = "sys_default"
	}

	transport, err := NewHttpTransport(config)
	if err != nil {
		return nil, errors.New(fmt.Errorf("failed to create transport: %w", err).Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	rid := GenerateRandomString(5)

	cl := &Client{
		Config: config,
		Ctx:    ctx,
		Cancel: cancel,
		HttpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		},
		AlgoID: "00000000-0000-0000-0000-000000000000",
		Log: log.New(
			os.Stdout,
			"["+rid+"][user:"+config.Username+" bind_device:"+config.BindInterface+"] ",
			log.LstdFlags|log.Lmsgprefix,
		),
		heartBeatTicker: time.NewTicker(time.Duration(math.MaxInt32)),
	}

	return cl, nil
}

func (c *Client) Start() {
	c.Log.Println("client start")
	defer wg.Done()
	defer c.Logout()

	if err := c.CheckNetwork(); err != nil {
		c.Log.Printf("Network check failed:%v", err)
	}

	ticker := time.NewTicker(time.Millisecond * time.Duration(c.Config.CheckInterval))
	defer ticker.Stop()

	for {
		select {
		case <-c.Ctx.Done():
			c.Log.Println("client context cancel")
			return
		case <-ticker.C:
			if err := c.CheckNetwork(); err != nil {
				c.Log.Printf("Network check failed:%v", err)
			}
		case <-c.heartBeatTicker.C:
			err := c.SendHeartbeat()
			if err != nil {
				c.Log.Printf("send heartbeat error: %v", err)
			} else {
				c.Log.Println("send heartbeat")
			}
		}
	}
}

func (c *Client) SendHeartbeat() error {
	stateXML, err := c.GenerateStateXML()
	if err != nil {
		return errors.New(err.Error())
	}

	decrypted, err := c.PostXML(c.KeepUrl, stateXML)

	var stateResp StateResponse
	if err := xml.Unmarshal(decrypted, &stateResp); err != nil {
		return errors.New(err.Error())
	}

	interval, err := strconv.Atoi(stateResp.Interval)
	if err != nil {
		return errors.New(err.Error())
	}

	c.heartBeatTicker.Reset(time.Duration(interval) * time.Second)
	return nil
}

func (c *Client) Logout() {
	request, _ := c.NewGetRequest("http://connect.rom.miui.com/generate_204")
	resp, _ := c.HttpClient.Do(request)
	if resp.StatusCode == http.StatusNoContent {
		stateXML, _ := c.GenerateStateXML()
		_, _ = c.PostXMLWithTimeout(c.TermUrl, stateXML)
		c.Log.Println("log out request sent")
	}
}

func (c *Client) CheckNetwork() error {
	request, err := c.NewGetRequest("http://connect.rom.miui.com/generate_204")
	if err != nil {
		return errors.New(err.Error())
	}

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return errors.New(err.Error())
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil

	case http.StatusFound:
		c.heartBeatTicker.Reset(time.Duration(math.MaxInt32))
		c.Log.Println("auth required")
		return c.HandleRedirect(resp)

	default:
		return errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}

func (c *Client) HandleRedirect(resp *http.Response) error {
	if err := c.Auth(resp.Header.Get("Location")); err != nil {
		c.Log.Printf("auth failed: %v", err)
		return nil
	}

	c.Log.Println("auth finished")
	return nil
}
