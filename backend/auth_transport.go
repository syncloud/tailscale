package main

import (
	"context"
	"net"
	"net/http"
)

type AuthTransport struct {
	authHost string
	socket   *http.Transport
	network  http.RoundTripper
}

func NewAuthTransport(authHost, socketPath string) *AuthTransport {
	return &AuthTransport{
		authHost: authHost,
		socket: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
			},
		},
		network: http.DefaultTransport,
	}
}

func (t *AuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Host != t.authHost {
		return t.network.RoundTrip(r)
	}
	viaSocket := r.Clone(r.Context())
	viaSocket.URL.Scheme = "http"
	viaSocket.Host = t.authHost
	viaSocket.Header.Set("X-Forwarded-Proto", "https")
	viaSocket.Header.Set("X-Forwarded-Host", t.authHost)
	return t.socket.RoundTrip(viaSocket)
}
