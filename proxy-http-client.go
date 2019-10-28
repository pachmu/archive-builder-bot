package main

import (
	"golang.org/x/net/proxy"
	"net/http"
)

func ProxyHttpClient(addr, username, password string) (*http.Client, error) {
	var auth *proxy.Auth = nil
	if username != "" || password != "" {
		auth = &proxy.Auth{
			User:     username,
			Password: password,
		}
	}
	dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}
	return httpClient, nil

}
