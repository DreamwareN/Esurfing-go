package client

import (
	"github.com/DreamwareN/Esurfing-go/cipher"
	"github.com/DreamwareN/Esurfing-go/client/utils"
	"github.com/google/uuid"
)

func (cl *Client) authorization(URL string) error {
	log := cl.Log
	cl.FirstRedirectURL = URL

	err := cl.getSchoolInfo()
	if err != nil {
		return err
	}

	log.Println("Domain: ", cl.Domain)
	log.Println("Area: ", cl.Area)
	log.Println("School Id: ", cl.SchoolID)
	log.Println("Index URL: ", cl.IndexURL)

	cl.ClientID = uuid.New()
	cl.Hostname = utils.GenerateRandomString(10)
	cl.MacAddress = utils.GenerateRandomMAC()

	//获取EConfig并设置TicketURL
	err = cl.getEConfig()
	if err != nil {
		return err
	}

	log.Println("Ticket URL: ", cl.TicketURL)

	err = cl.getUserAndAcIP()
	if err != nil {
		return err
	}

	log.Println("User IP: ", cl.UserIP)
	log.Println("AcIP: ", cl.AcIP)

	//get algo id
	err = cl.getAlgoId()
	if err != nil {
		return err
	}

	cl.cipher = cipher.NewCipher(cl.AlgoID)

	log.Println("Algo ID:", cl.AlgoID)
	log.Println("Key:", cl.Key)

	//get ticket
	err = cl.getTicket()
	if err != nil {
		return err
	}

	log.Println("Ticket: ", cl.Ticket)

	//login
	err = cl.login()
	if err != nil {
		return err
	}

	log.Println("Keep URL:", cl.KeepURL)
	log.Println("Term URL:", cl.TermURL)

	return nil
}
