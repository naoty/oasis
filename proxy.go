package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

type Proxy struct {
	URL        *url.URL
	BackendURL *url.URL
	Repository *Repository
}

func NewProxy(host, backendHost, repositoryURLString string) *Proxy {
	hostURL, err := url.Parse(normalizeURLString(host))
	if err != nil {
		log.Fatalf("URL parse error: %s", err)
	}

	backendURL, err := url.Parse(normalizeURLString(backendHost))
	if err != nil {
		log.Fatalf("URL parse error: %s", err)
	}

	repository := NewRepository(repositoryURLString, "master")

	return &Proxy{URL: hostURL, BackendURL: backendURL, Repository: repository}
}

func (p Proxy) Start() error {
	index := LoadIndex()
	director := func(req *http.Request) {
		subdomain := parseSubdomain(req.Host)

		port, err := index.LookupPort(p.Repository.RemoteURL.String(), subdomain)
		if err != nil {
			// TODO: Inspect the port of the docker host for a container
			p.Repository.Checkout(subdomain)
			port = ""
			// revision := p.Repository.Checkout(subdomain)
			// revision.BuildAndRun()
			// port = revision.InspectPort()
		}

		req.URL = rewriteURL(req.URL, p.BackendURL)
		req.URL.Host = fmt.Sprintf("%s:%s", req.URL.Host, port)

		log.Printf("Redirect to: %s\n", req.URL.String())
	}

	reverseProxy := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    p.URL.Host,
		Handler: reverseProxy,
	}

	fmt.Printf("Listening: %s\n", p.URL.Host)
	return server.ListenAndServe()
}

var schemePattern = regexp.MustCompile("^[^:]+://")

func normalizeURLString(urlString string) string {
	if schemePattern.MatchString(urlString) {
		return urlString
	} else {
		return "http://" + urlString
	}
}

func parseSubdomain(host string) string {
	labels := strings.Split(host, ".")
	return labels[0]
}

func rewriteURL(originalURL, backendURL *url.URL) (newURL *url.URL) {
	newURL = new(url.URL)

	newURL.Scheme = backendURL.Scheme
	newURL.Host = backendURL.Host
	newURL.Path = singleJoiningSlash(backendURL.Path, originalURL.Path)

	if backendURL.RawQuery == "" || originalURL.RawQuery == "" {
		newURL.RawQuery = backendURL.RawQuery + originalURL.RawQuery
	} else {
		newURL.RawQuery = backendURL.RawQuery + "&" + originalURL.RawQuery
	}
	return newURL
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
