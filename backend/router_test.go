package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRouter(t *testing.T) (*httptest.Server, *SessionCodec, *FakeExecutor) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("headscale " + r.URL.Path))
	}))
	t.Cleanup(upstream.Close)
	target, err := url.Parse(upstream.URL)
	require.NoError(t, err)

	webDir := t.TempDir()
	require.NoError(t, os.WriteFile(path.Join(webDir, "index.html"), []byte("<div id=\"app\"></div>"), 0644))
	require.NoError(t, os.MkdirAll(path.Join(webDir, "assets"), 0755))
	require.NoError(t, os.WriteFile(path.Join(webDir, "assets", "app.js"), []byte("js"), 0644))

	codec := NewSessionCodec([]byte("key"))
	secret := func() (string, error) { return "secret", nil }
	auth := NewOidc("https://auth.device", "tailscale", secret, "https://tailscale.device"+CallbackPath, http.DefaultClient, codec)
	devices, executor := fakeDevices()
	server := httptest.NewServer(NewRouter("https://tailscale.device", webDir, auth, devices, NewHeadscaleProxy(target)))
	t.Cleanup(server.Close)
	return server, codec, executor
}

func request(t *testing.T, method, url, cookie string) (int, string) {
	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookie, Value: cookie})
	}
	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response.StatusCode, string(body)
}

func TestRouterServesSpaAndAssets(t *testing.T) {
	server, _, _ := testRouter(t)
	status, body := request(t, "GET", server.URL+"/", "")
	assert.Equal(t, 200, status)
	assert.Contains(t, body, "id=\"app\"")
	status, body = request(t, "GET", server.URL+"/app/assets/app.js", "")
	assert.Equal(t, 200, status)
	assert.Equal(t, "js", body)
}

func TestRouterSendsEverythingElseToHeadscale(t *testing.T) {
	server, _, _ := testRouter(t)
	for _, p := range []string{"/health", "/api/v1/node", "/register/abc", "/oidc/callback", "/derp/probe"} {
		status, body := request(t, "GET", server.URL+p, "")
		assert.Equal(t, 200, status, p)
		assert.Equal(t, "headscale "+p, body)
	}
}

func TestRouterApiNeedsSession(t *testing.T) {
	server, _, _ := testRouter(t)
	status, _ := request(t, "GET", server.URL+"/app/api/devices", "")
	assert.Equal(t, 401, status)
	status, _ = request(t, "GET", server.URL+"/app/api/devices", "forged.value")
	assert.Equal(t, 401, status)
}

func TestRouterLoginBeforeDiscoveryIsUnavailable(t *testing.T) {
	server, _, _ := testRouter(t)
	status, _ := request(t, "GET", server.URL+LoginPath, "")
	assert.Equal(t, 503, status)
}

func TestRouterListsAndRemovesDevices(t *testing.T) {
	server, codec, executor := testRouter(t)
	cookie, err := codec.Encode(Session{User: "boris", Expires: time.Now().Add(time.Hour).Unix()})
	require.NoError(t, err)

	status, body := request(t, "GET", server.URL+"/app/api/me", cookie)
	assert.Equal(t, 200, status)
	assert.JSONEq(t, `{"user":"boris","url":"https://tailscale.device"}`, body)

	status, body = request(t, "GET", server.URL+"/app/api/devices", cookie)
	assert.Equal(t, 200, status)
	assert.Contains(t, body, `"name":"phone"`)

	status, _ = request(t, "DELETE", server.URL+"/app/api/devices/1", cookie)
	assert.Equal(t, 403, status)
	status, _ = request(t, "DELETE", server.URL+"/app/api/devices/2", cookie)
	assert.Equal(t, 204, status)
	assert.Contains(t, executor.calls, "hs nodes delete --identifier 2 --force")
}
