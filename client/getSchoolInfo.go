package client

import (
	"github.com/DreamwareN/Esurfing-go/errs"
)

func (cl *Client) getSchoolInfo() error {
	if cl.FirstRedirectURL == "" {
		return errs.New("missing redirect URL")
	}

	request, err := cl.GenerateGetRequest(cl.FirstRedirectURL)
	if err != nil {
		return errs.New(err.Error())
	}

	response, err := cl.HttpClient.Do(request)
	if err != nil {
		return errs.New(err.Error())
	}

	if response.Header.Get("domain") != "" && response.Header.Get("area") != "" &&
		response.Header.Get("schoolid") != "" && response.Header.Get("Location") != "" {
		cl.Domain = response.Header.Get("domain")
		cl.Area = response.Header.Get("area")
		cl.SchoolID = response.Header.Get("schoolid")
		cl.IndexURL = response.Header.Get("Location")
	} else {
		return errs.New("missing school info")
	}

	if response.StatusCode != 302 {
		return errs.New("invalid process of authorization at stage 2")
	}

	cl.SecondRedirectURL = response.Header.Get("location")
	return nil
}
