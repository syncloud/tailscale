package main

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordsPointDeviceAndAppsAtNode(t *testing.T) {
	snapDir := t.TempDir()
	for _, dir := range []string{"nextcloud", "immich", "platform", "tailscale", "core22", "snapd", "bin"} {
		require.NoError(t, os.Mkdir(path.Join(snapDir, dir), 0755))
	}
	require.NoError(t, os.WriteFile(path.Join(snapDir, "README"), nil, 0644))
	file := path.Join(t.TempDir(), "records.json")

	err := NewRecords(file, "device.syncloud.it", snapDir).Sync([]string{"100.64.0.1", "fd7a:115c:a1e0::1"})
	require.NoError(t, err)

	var records []Record
	content, err := os.ReadFile(file)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(content, &records))
	assert.Equal(t, []Record{
		{Name: "device.syncloud.it", Type: "A", Value: "100.64.0.1"},
		{Name: "device.syncloud.it", Type: "AAAA", Value: "fd7a:115c:a1e0::1"},
		{Name: "auth.device.syncloud.it", Type: "A", Value: "100.64.0.1"},
		{Name: "auth.device.syncloud.it", Type: "AAAA", Value: "fd7a:115c:a1e0::1"},
		{Name: "immich.device.syncloud.it", Type: "A", Value: "100.64.0.1"},
		{Name: "immich.device.syncloud.it", Type: "AAAA", Value: "fd7a:115c:a1e0::1"},
		{Name: "nextcloud.device.syncloud.it", Type: "A", Value: "100.64.0.1"},
		{Name: "nextcloud.device.syncloud.it", Type: "AAAA", Value: "fd7a:115c:a1e0::1"},
	}, records)
}

func TestRecordsNotRewrittenWhenUnchanged(t *testing.T) {
	snapDir := t.TempDir()
	file := path.Join(t.TempDir(), "records.json")
	records := NewRecords(file, "device.syncloud.it", snapDir)
	require.NoError(t, records.Sync([]string{"100.64.0.1"}))
	before, err := os.Stat(file)
	require.NoError(t, err)

	require.NoError(t, records.Sync([]string{"100.64.0.1"}))
	after, err := os.Stat(file)
	require.NoError(t, err)
	assert.Equal(t, before.ModTime(), after.ModTime())
}
