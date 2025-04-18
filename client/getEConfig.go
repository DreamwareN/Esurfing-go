package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/errs"
	"github.com/DreamwareN/Esurfing-go/utils"
	"io"
)

func (cl *Client) GetEConfig() error {
	if cl.FirstRedirectURL == "" {
		return errs.New("missing redirect URL")
	}

	request, err := cl.GenerateGetRequest(cl.IndexURL)
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
