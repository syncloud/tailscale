package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyServesPageAtRoot(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("headscale " + r.URL.Path))
	}))
	defer upstream.Close()
	target, err := url.Parse(upstream.URL)
	require.NoError(t, err)
	server := httptest.NewServer(NewProxy("", target, NewPage("https://tailscale.device.syncloud.it")))
	defer server.Close()

	root := get(t, server.URL+"/")
	assert.Contains(t, root, "https://tailscale.device.syncloud.it")

	health := get(t, server.URL+"/health")
	assert.Equal(t, "headscale /health", health)
}

func get(t *testing.T, url string) string {
	response, err := http.Get(url)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return string(body)
}
