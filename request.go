package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

func (c *Client) NewGetRequest(url string) (request *http.Request, err error) {
	req, err := http.NewRequestWithContext(c.Ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", c.ClientID.String())
	req.Header.Set("Connection", "keep-alive")
	if c.SchoolID != "" {
		req.Header.Set("CDC-SchoolId", c.SchoolID)
	}
	if c.Domain != "" {
		req.Header.Set("CDC-Domain", c.Domain)
	}
	if c.Area != "" {
		req.Header.Set("CDC-Area", c.Area)
	}

	return req, nil
}

func (c *Client) NewPostRequest(url string, data []byte) (request *http.Request, err error) {
	md5Hex := md5.Sum(data)

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", c.ClientID.String())
	req.Header.Set("CDC-Checksum", hex.EncodeToString(md5Hex[:]))
	req.Header.Set("Algo-ID", c.AlgoID)
	return req, nil
}

func (c *Client) NewPostRequestWithCustomCtx(ctx context.Context, url string, data []byte) (request *http.Request, err error) {
	md5Hex := md5.Sum(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgentAndroid)
	req.Header.Set("Accept", "text/html,text/xml,application/xhtml+xml,application/x-javascript,*/*")
	req.Header.Set("Client-ID", c.ClientID.String())
	req.Header.Set("CDC-Checksum", hex.EncodeToString(md5Hex[:]))
	req.Header.Set("Algo-ID", c.AlgoID)
	return req, nil
}
