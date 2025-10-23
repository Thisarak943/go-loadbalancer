package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleserver struct {
	addr  string
	proxy httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleserver {
	serverUrl, err := url.Parse(addr)
	handleErr(err)
	return &simpleserver{
		addr:  addr,
		proxy: *httputil.NewSingleHostReverseProxy(serverUrl),
	}

}
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

