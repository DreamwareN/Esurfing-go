package client

import (
	"fmt"
	"github.com/DreamwareN/Esurfing-go/errs"
	"net/http"
	"time"
)

func (cl *Client) CheckNetworkStatus() error {
	request, err := cl.GenerateGetRequest("http://connect.rom.miui.com/generate_204")
	if err != nil {
		return errs.New(err.Error())
	}

	resp, err := cl.HttpClient.Do(request)
	if err != nil {
		return errs.New(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		cl.IsRunning.Store(1)
		if cl.IsConnectedAtFirst == 0 {
			cl.IsConnectedAtFirst = 1
			cl.Log.Println("The network has been connected")
		}
		return nil

	case http.StatusFound:
		cl.Log.Println("Authorization required")
		return cl.HandleRedirect(resp)

	default:
		return errs.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}

func (cl *Client) HandleRedirect(resp *http.Response) error {
	redirectURL := resp.Header.Get("Location")
	if redirectURL == "" {
		return errs.New("missing Location header in redirect")
	}

	retryCount := 0
	for retryCount <= cl.Conf.MaxRetries {
		select {
		case <-cl.Ctx.Done():
			return errs.New("context canceled")
		default:
		}

		if err := cl.Authorization(redirectURL); err != nil {
			retryCount++
			cl.Log.Printf("Authorization attempt %d failed: %v", retryCount, err)
			select {
			case <-time.After(time.Duration(cl.Conf.RetryDelayMS) * time.Millisecond):
			case <-cl.Ctx.Done():
				return errs.New("context canceled")
			}
			continue
		}

		cl.Log.Println("Authorization succeeded")
		cl.IsLogin.Store(1)
		cl.IsRunning.Store(1)
		go cl.MaintainSession()
		return nil
	}

	cl.Log.Println("Reached max retry count: ", retryCount)
	return ErrMaxReTryReach
}
