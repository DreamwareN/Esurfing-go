package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/errs"
	"strconv"
	"time"
)

func (cl *Client) sendHeartbeat() error {
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
