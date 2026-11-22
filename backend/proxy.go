package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Proxy struct {
	socket string
	page   http.Handler
	target *httputil.ReverseProxy
}

func NewProxy(socket string, upstream *url.URL, page http.Handler) *Proxy {
	target := httputil.NewSingleHostReverseProxy(upstream)
	target.FlushInterval = -1
	return &Proxy{socket: socket, page: page, target: target}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		p.page.ServeHTTP(w, r)
		return
	}
	p.target.ServeHTTP(w, r)
}

func (p *Proxy) Start() error {
	_ = os.Remove(p.socket)
	listener, err := net.Listen("unix", p.socket)
	if err != nil {
		return err
	}
	if err := os.Chmod(p.socket, 0666); err != nil {
		return err
	}
	return http.Serve(listener, p)
}
