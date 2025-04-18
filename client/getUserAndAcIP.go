package client

import (
	"github.com/DreamwareN/Esurfing-go/errs"
	"net/url"
)

func (cl *Client) GetUserAndAcIP() error {
	URLParsed, err := url.Parse(cl.TicketURL)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.UserIP = URLParsed.Query().Get("wlanuserip")
	cl.AcIP = URLParsed.Query().Get("wlanacip")

	if cl.UserIP == "" || cl.AcIP == "" {
		return errs.New("missing user ip or ac ip")
	}

	return nil
}
