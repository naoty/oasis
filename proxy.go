package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	Host        string
	BackendHost string
	Repository  string
}

func (p Proxy) Start() error {
	backendURL, err := url.Parse(p.BackendHost)
	if err != nil {
		log.Fatalf("While parsing %s: %s", p.BackendHost, err)
	}

	index := LoadIndex()
	backendQuery := backendURL.RawQuery
	director := func(req *http.Request) {
		subdomain := parseSubdomain(req.Host)
		port, err := index.LookupPort(p.Repository, subdomain)

		if err != nil {
			workspace := NewWorkspace(p.Repository, subdomain)
			workspace.Run()
			port = workspace.InspectPort()
		}

		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host + ":" + port
		req.URL.Path = singleJoiningSlash(backendURL.Path, req.URL.Path)

		if backendQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = backendQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = backendQuery + "&" + req.URL.RawQuery
		}

		log.Printf("Redirect to: %s\n", req.URL.String())
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    p.Host,
		Handler: reverseProxy,
	}

	fmt.Printf("Listening: %s\n", p.Host)
	return server.ListenAndServe()
}

func parseSubdomain(host string) string {
	labels := strings.Split(host, ".")
	return labels[0]
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasSuffix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
