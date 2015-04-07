package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
)

type Proxy struct {
	URL              *url.URL
	ContainerHostURL *url.URL
	Workspace        *Workspace
	Index            *Index
}

func NewProxy(proxyURL, containerHostURL, repositoryURL *url.URL) *Proxy {
	index := NewIndex()
	workspace := NewWorkspace(repositoryURL, containerHostURL, index)
	return &Proxy{
		URL:              proxyURL,
		ContainerHostURL: containerHostURL,
		Workspace:        workspace,
		Index:            index,
	}
}

func (proxy *Proxy) Start() error {
	director := proxy.newDirector()
	reverseProxy := &httputil.ReverseProxy{Director: director}
	server := http.Server{Addr: proxy.URL.Host, Handler: reverseProxy}

	logrus.WithFields(logrus.Fields{
		"url": proxy.URL.String(),
	}).Info("Start a proxy")

	return server.ListenAndServe()
}

func (proxy *Proxy) newDirector() func(request *http.Request) {
	return func(request *http.Request) {
		subdomain := proxy.parseSubdomain(request.Host)
		port, err := proxy.Workspace.LookupPort(subdomain)

		if err != nil {
			port = proxy.Workspace.Setup(subdomain)
		}

		targetURL := proxy.rewriteURL(request.URL, port)

		logrus.WithFields(logrus.Fields{
			"target": targetURL.String(),
		}).Info("Redirect a request")

		request.URL = targetURL
	}
}

func (proxy *Proxy) parseSubdomain(host string) string {
	labels := strings.Split(host, ".")
	return labels[0]
}

func (proxy *Proxy) rewriteURL(originalURL *url.URL, port string) *url.URL {
	return &url.URL{
		Scheme:   proxy.rewriteScheme(originalURL.Scheme),
		Host:     proxy.rewritePort(proxy.ContainerHostURL.Host, port),
		Path:     originalURL.Path,
		RawQuery: originalURL.RawQuery,
	}
}

func (proxy *Proxy) rewriteScheme(scheme string) string {
	if scheme == "" {
		return "http"
	} else {
		return scheme
	}
}

func (proxy *Proxy) rewritePort(host, port string) string {
	elements := strings.Split(host, ":")
	return fmt.Sprintf("%s:%s", elements[0], port)
}
