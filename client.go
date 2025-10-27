package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Config     *Config
	Log        *log.Logger
	HttpClient *http.Client
	Ctx        context.Context
	Wg         *sync.WaitGroup
	cipher     Cipher

	//mark whether the network is connected
	isAuthenticated bool

	UserIP            string
	AcIP              string
	Domain            string
	Area              string
	SchoolID          string
	ClientID          uuid.UUID
	Hostname          string
	MacAddress        string
	Ticket            string
	AlgoID            string
	HeartbeatInterval time.Duration

	IndexUrl    string
	TicketUrl   string
	AuthUrl     string
	KeepUrl     string
	TermUrl     string
	RedirectUrl string
}

func (c *Client) Start() {
	c.Log.Println("client start")
	defer c.Wg.Done()

	networkCheckTicker := time.NewTicker(time.Millisecond * time.Duration(c.Config.CheckInterval))
	defer networkCheckTicker.Stop()

	for {
		select {
		case <-networkCheckTicker.C:
			//remove go routine due to potential memory and cpu overload
			if err := c.CheckNetwork(); err != nil {
				c.Log.Printf("Network check failed:%v", err)
			}
		case <-c.Ctx.Done():
			c.Log.Println("client stop")
			if c.isAuthenticated {
				c.Logout()
				c.Log.Println("log out request sent")
			}
			return
		}
	}

}

func (c *Client) Logout() {
	stateXML, _ := c.GenerateStateXML()
	_, _ = c.PostXMLWithTimeout(c.TermUrl, stateXML)
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
		//network is connected
		return nil

	case http.StatusFound:
		//swap the status and do authenticate
		c.isAuthenticated = false
		c.Log.Println("auth required")
		return c.HandleRedirect(resp)

	default:
		return errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}

func (c *Client) HandleRedirect(resp *http.Response) error {
	for {
		if err := c.Auth(resp.Header.Get("Location")); err != nil {
			c.Log.Printf("auth failed: %v", err)
			select {
			case <-time.After(time.Duration(c.Config.RetryInterval) * time.Millisecond):
			case <-c.Ctx.Done():
				return errors.New("context canceled")
			}
			continue
		}

		c.Log.Println("auth finished")
		c.isAuthenticated = true
		go c.MaintainSession()
		return nil
	}
}
