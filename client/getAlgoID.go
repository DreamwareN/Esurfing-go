package client

import (
	"github.com/DreamwareN/Esurfing-go/client/utils"
	"github.com/DreamwareN/Esurfing-go/errs"
	"io"
)

func (cl *Client) getAlgoId() error {
	request, err := cl.GeneratePostRequest(cl.TicketURL, []byte(cl.AlgoID))
	if err != nil {
		return errs.New(err.Error())
	}

	response, err := cl.HttpClient.Do(request)
	if err != nil {
		return errs.New(err.Error())
	}

	defer response.Body.Close()

	algoIdData, err := io.ReadAll(response.Body)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.AlgoID, cl.Key, err = utils.DecodeAlgoID(algoIdData)
	if err != nil {
		return errs.New(err.Error())
	}
	return nil
}
