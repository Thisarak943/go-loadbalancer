package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Addr() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

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

type LoadBalancer struct {
	port              string
	roundRobinCounter int
	servers           []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:              port,
		roundRobinCounter: 0,
		servers:           servers,
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	// implementation will be added later
	return nil
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	// implementation will be added later
}
