package probing

import (
	"crypto/tls"
	"net/http"
	"time"
)

func NewHttpClient(timeoutSecond, tlsHandsakeTimeout int, ignoreTLS bool) *http.Client {
	netTransport := &http.Transport{
		TLSHandshakeTimeout: time.Duration(tlsHandsakeTimeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: ignoreTLS},
	}

	ret := &http.Client{
		Timeout:   time.Duration(timeoutSecond) * time.Second,
		Transport: netTransport,
	}

	return ret
}
