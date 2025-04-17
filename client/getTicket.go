package client

import (
	"encoding/xml"
	"github.com/DreamwareN/Esurfing-go/errs"
)

func (cl *Client) getTicket() error {
	getTicketXML, err := cl.GenerateGetTicketXML()
	if err != nil {
		return errs.New(err.Error())
	}

	ticketData, err := cl.PostXML(cl.TicketURL, getTicketXML)
	if err != nil {
		return errs.New(err.Error())
	}

	ticketXML := &TicketResponse{}

	err = xml.Unmarshal(ticketData, ticketXML)
	if err != nil {
		return errs.New(err.Error())
	}

	cl.Ticket = ticketXML.Ticket
	return nil
}
