package client

import (
	"encoding/xml"
	"io"
	"time"
)

func (cl *Client) GenerateGetTicketXML() ([]byte, error) {
	tr := TicketRequest{
		UserAgent: UserAgent,
		ClientID:  cl.ClientID.String(),
		LocalTime: time.Now().Format(time.DateTime),
		HostName:  cl.Hostname,
		Ipv4:      cl.UserIP,
		Mac:       cl.MacAddress,
		Ostag:     cl.Hostname,
		Gwip:      cl.AcIP,
	}
	out, err := xml.Marshal(tr)
	if err != nil {
		return nil, err
	}
	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), out...), nil
}

func (cl *Client) GenerateStateXML() ([]byte, error) {
	s := &State{
		UserAgent: UserAgent,
		ClientID:  cl.ClientID.String(),
		LocalTime: time.Now().Format(time.DateTime),
		HostName:  cl.Hostname,
		Ipv4:      cl.UserIP,
		Ticket:    cl.Ticket,
		Mac:       cl.MacAddress,
		Ostag:     cl.Hostname,
	}
	bytes, err := xml.Marshal(s)
	if err != nil {
		return nil, err
	}

	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), bytes...), nil
}

func (cl *Client) GenerateLoginXML() ([]byte, error) {
	lr := &LoginRequest{
		UserAgent: UserAgent,
		ClientID:  cl.ClientID.String(),
		Ticket:    cl.Ticket,
		LocalTime: time.Now().Format(time.DateTime),
		Userid:    cl.Conf.AuthUsername,
		Passwd:    cl.Conf.AuthPassword,
	}

	bytes, err := xml.Marshal(lr)
	if err != nil {
		return nil, err
	}

	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"), bytes...), nil
}

func (cl *Client) PostXML(url string, data []byte) ([]byte, error) {
	encXML, err := cl.cipher.Encrypt(data)
	if err != nil {
		return nil, err
	}

	req, err := cl.GeneratePostRequest(url, encXML)
	if err != nil {
		return nil, err
	}

	response, err := cl.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return cl.cipher.Decrypt(data)
}

func (cl *Client) PostXMLWithoutCtx(url string, data []byte) ([]byte, error) {
	encXML, err := cl.cipher.Encrypt(data)
	if err != nil {
		return nil, err
	}

	req, err := cl.GeneratePostRequestWithoutCtx(url, encXML)
	if err != nil {
		return nil, err
	}

	response, err := cl.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return cl.cipher.Decrypt(data)
}
