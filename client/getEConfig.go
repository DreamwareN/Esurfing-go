package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/client/utils"
	"github.com/DreamwareN/Esurfing-go/errs"
	"io"
)

func (cl *Client) getEConfig() error {
	if cl.FirstRedirectURL == "" {
		return errs.New("missing redirect URL")
	}

	request, err := cl.GenerateGetRequest(cl.SecondRedirectURL)
	if err != nil {
		return errs.New(err.Error())
	}

	response, err := cl.HttpClient.Do(request)
	if err != nil {
		return errs.New(err.Error())
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return errs.New(err.Error())
	}

	eConfigData, err := utils.ParseEConfig(data)
	if err != nil {
		return errs.New(err.Error())
	}

	eConfig := &EConfig{}

	err = xml.Unmarshal(eConfigData, eConfig)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.TicketURL = eConfig.TicketURL
	cl.AuthURL = eConfig.AuthURL

	return nil
}
