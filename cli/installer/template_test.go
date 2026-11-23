package installer

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syncloud/golib/config"
	"gopkg.in/yaml.v3"
)

func TestTemplatesRender(t *testing.T) {
	out := t.TempDir()
	variables := Variables{
		App:           App,
		AppDir:        "/snap/tailscale/current",
		CommonDir:     "/var/snap/tailscale/common",
		DataDir:       "/var/snap/tailscale/current",
		Domain:        "device.syncloud.it",
		AppDomain:     "tailscale.device.syncloud.it",
		Hostname:      "device",
		Url:           "https://tailscale.device.syncloud.it",
		AuthUrl:       "https://auth.device.syncloud.it",
		AuthSocket:    "/var/snap/platform/current/authelia.socket",
		StunHost:      "stun.syncloud.it",
		StunPort:      StunPort,
		HeadscalePort: HeadscalePort,
		MetricsPort:   MetricsPort,
		GrpcPort:      GrpcPort,
		DerpStunPort:  DerpStunPort,
	}
	require.NoError(t, config.Generate("../../config", out, variables))

	headscale := map[string]any{}
	content, err := os.ReadFile(path.Join(out, "headscale.yaml"))
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(content, &headscale))
	assert.Equal(t, "https://tailscale.device.syncloud.it", headscale["server_url"])
	split := headscale["dns"].(map[string]any)["nameservers"].(map[string]any)["split"].(map[string]any)
	assert.Contains(t, split, "device.syncloud.it")

	derp := map[string]any{}
	content, err = os.ReadFile(path.Join(out, "derp.yaml"))
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(content, &derp))
	nodes := derp["regions"].(map[any]any)[900].(map[string]any)["nodes"].([]any)
	assert.Equal(t, "stun.syncloud.it", nodes[0].(map[string]any)["hostname"])
	assert.Equal(t, "tailscale.device.syncloud.it", nodes[1].(map[string]any)["hostname"])

	backend := map[string]any{}
	content, err = os.ReadFile(path.Join(out, "backend.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(content, &backend))
	assert.Equal(t, "/snap/tailscale/current/web", backend["web_dir"])
	assert.Equal(t, "/var/snap/platform/current/authelia.socket", backend["auth_socket"])
	assert.Equal(t, "https://auth.device.syncloud.it", backend["auth_url"])
}
