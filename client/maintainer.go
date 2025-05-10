package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/errs"
	"strconv"
	"time"
)

func (cl *Client) MaintainSession() {
	cl.IsLogin.Store(true)
	cl.IsRunning.Store(true)
	cl.Log.Println("Maintaining session")
	for cl.IsRunning.Load() {
		cl.Log.Println("Next Heartbeat: ", time.Now().Add(cl.HeartbeatInterval).Format(time.DateTime))
		select {
		case <-time.After(cl.HeartbeatInterval):
			if err := cl.SendHeartbeat(); err != nil {
				cl.Log.Printf("Heartbeat failed: %v", err)
				cl.IsRunning.Store(false)
				return
			}
			cl.Log.Println("Heartbeat Sent")
		case <-cl.Ctx.Done():
			return
		}
	}
}

func (cl *Client) SendHeartbeat() error {
	stateXML, err := cl.GenerateStateXML()
	if err != nil {
		return errs.New(err.Error())
	}

	decrypted, err := cl.PostXML(cl.KeepURL, stateXML)

	var stateResp StateResponse
	if err := xml.Unmarshal(decrypted, &stateResp); err != nil {
		return errs.New(err.Error())
	}

	interval, err := strconv.Atoi(stateResp.Interval)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.HeartbeatInterval = time.Duration(interval) * time.Second
	return nil
}
