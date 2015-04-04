package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Host        string
	BackendHost string
}

func (p Proxy) Start() error {
	backendUrl, err := url.Parse(p.BackendHost)
	if err != nil {
		log.Fatalf("While parsing %s: %s", p.BackendHost, err)
	}

	proxyHandler := httputil.NewSingleHostReverseProxy(backendUrl)
	server := http.Server{
		Addr:    p.Host,
		Handler: proxyHandler,
	}
	fmt.Printf("Listening: %s\n", p.Host)
	return server.ListenAndServe()
}
