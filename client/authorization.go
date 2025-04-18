package client

import (
	"github.com/DreamwareN/Esurfing-go/cipher"
	"github.com/DreamwareN/Esurfing-go/errs"
	"github.com/DreamwareN/Esurfing-go/utils"
	"github.com/google/uuid"
)

func (cl *Client) Authorization(URL string) error {
	log := cl.Log
	cl.FirstRedirectURL = URL

	err := cl.GetSchoolInfo()
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
	err = cl.GetEConfig()
	if err != nil {
		return err
	}

	log.Println("Ticket URL: ", cl.TicketURL)

	err = cl.GetUserAndAcIP()
	if err != nil {
		return err
	}

	log.Println("User IP: ", cl.UserIP)
	log.Println("Ac IP: ", cl.AcIP)

	//get algo id
	err = cl.GetAlgoId()
	if err != nil {
		return err
	}

	cl.cipher = cipher.NewCipher(cl.AlgoID)
	if cl.cipher == nil {
		return errs.New("Unknown AlgoID: " + cl.AlgoID)
	}

	log.Println("Algo ID:", cl.AlgoID)

	//get ticket
	err = cl.GetTicket()
	if err != nil {
		return err
	}

	log.Println("Ticket: ", cl.Ticket)

	//login
	err = cl.Login()
	if err != nil {
		return err
	}

	log.Println("Keep URL:", cl.KeepURL)
	log.Println("Term URL:", cl.TermURL)

	return nil
}
