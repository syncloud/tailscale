package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewHeadscaleProxy(upstream *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	proxy.FlushInterval = -1
	return proxy
}
