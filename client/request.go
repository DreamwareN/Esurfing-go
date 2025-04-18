package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

func (cl *Client) GenerateGetRequest(url string) (request *http.Request, err error) {
	req, err := http.NewRequestWithContext(cl.Ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", cl.ClientID.String())
	req.Header.Set("Connection", "keep-alive")
	if cl.SchoolID != "" {
		req.Header.Set("CDC-SchoolId", cl.SchoolID)
	}
	if cl.Domain != "" {
		req.Header.Set("CDC-Domain", cl.Domain)
	}
	if cl.Area != "" {
		req.Header.Set("CDC-Area", cl.Area)
	}

	return req, nil
}

func (cl *Client) GeneratePostRequest(url string, data []byte) (request *http.Request, err error) {
	md5Hex := md5.Sum(data)

	req, err := http.NewRequestWithContext(cl.Ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", cl.ClientID.String())
	req.Header.Set("CDC-Checksum", hex.EncodeToString(md5Hex[:]))
	req.Header.Set("Algo-ID", cl.AlgoID)
	return req, nil
}

func (cl *Client) GeneratePostRequestWithSpecCtx(ctx context.Context, url string, data []byte) (request *http.Request, err error) {
	md5Hex := md5.Sum(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", cl.ClientID.String())
	req.Header.Set("CDC-Checksum", hex.EncodeToString(md5Hex[:]))
	req.Header.Set("Algo-ID", cl.AlgoID)
	return req, nil
}
