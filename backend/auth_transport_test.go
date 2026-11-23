package main

import (
	"io"
	"net"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHostGoesOverSocketWithForwardedHeaders(t *testing.T) {
	socket := path.Join(t.TempDir(), "authelia.socket")
	listener, err := net.Listen("unix", socket)
	require.NoError(t, err)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Host + " " + r.Header.Get("X-Forwarded-Proto") + " " + r.Header.Get("X-Forwarded-Host") + " " + r.URL.Path))
	})}
	go server.Serve(listener)
	t.Cleanup(func() { _ = server.Close() })

	client := &http.Client{Transport: NewAuthTransport("auth.device.syncloud.it", socket)}
	response, err := client.Get("https://auth.device.syncloud.it/.well-known/openid-configuration")
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "auth.device.syncloud.it https auth.device.syncloud.it /.well-known/openid-configuration", string(body))
}
