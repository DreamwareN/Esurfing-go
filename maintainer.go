package main

import (
	"encoding/xml"
	"errors"
	"strconv"
	"time"
)

func (c *Client) MaintainSession() {
	c.Log.Println("Maintaining session")
	for c.isAuthenticated {
		c.Log.Println("next heartbeat: ", time.Now().Add(c.HeartbeatInterval).Format(time.DateTime))
		select {
		case <-time.After(c.HeartbeatInterval):
			if err := c.SendHeartbeat(); err != nil {
				c.Log.Printf("send heartbeat error: %v", err)
				return
			}
		case <-c.Ctx.Done():
			return
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

	c.HeartbeatInterval = time.Duration(interval) * time.Second
	return nil
}
