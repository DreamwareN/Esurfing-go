package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/errs"
	"strconv"
	"time"
)

func (cl *Client) login() error {
	loginXML, err := cl.GenerateLoginXML()
	if err != nil {
		return errs.New(err.Error())
	}

	responseData, err := cl.PostXML(cl.AuthURL, loginXML)
	if err != nil {
		return errs.New(err.Error())
	}

	loginResponseXML := &LoginResponse{}
	err = xml.Unmarshal(responseData, loginResponseXML)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.KeepURL = loginResponseXML.KeepURL
	cl.TermURL = loginResponseXML.TermURL

	keepRetrySec, err := strconv.Atoi(loginResponseXML.KeepRetry)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.HeartbeatInterval = time.Second * time.Duration(keepRetrySec)
	return nil
}
