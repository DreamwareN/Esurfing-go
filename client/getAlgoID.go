package client

import (
	"github.com/DreamwareN/Esurfing-go/errs"
	"github.com/DreamwareN/Esurfing-go/utils"
	"io"
)

func (cl *Client) GetAlgoId() error {
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

	cl.AlgoID, _, err = utils.DecodeAlgoID(algoIdData)
	if err != nil {
		return errs.New(err.Error())
	}
	return nil
}
