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

// ✅ Fixed: method name must match interface (Addr not Address)
func (s *simpleserver) Addr() string {
	return s.addr
}

func (s *simpleserver) IsAlive() bool {
	return true
}

func (s *simpleserver) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCounter%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCounter++
		server = lb.servers[lb.roundRobinCounter%len(lb.servers)]
	}
	lb.roundRobinCounter++
	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n", targetServer.Addr())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", servers)

	// ✅ Fixed: Correct handler registration
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	})

	fmt.Printf("Serving requests at 'localhost:%s'\n", lb.port)
	handleErr(http.ListenAndServe(":"+lb.port, nil))
}
