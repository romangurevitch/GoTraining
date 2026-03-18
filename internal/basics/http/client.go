package http

import (
	"net/http"
	"time"
)

func HowToCreateAnHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                  nil,
			DialContext:            nil,
			DialTLSContext:         nil,
			TLSClientConfig:        nil,
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar:     nil,
		Timeout: time.Duration(10) * time.Second,
	}
	// client.Do()
	return client
}
