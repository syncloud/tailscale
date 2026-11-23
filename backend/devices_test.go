package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nodesJson = `[
 {"id":1,"name":"device","given_name":"device","user":{"id":1,"name":"device"},"ip_addresses":["100.64.0.1"],"online":true,"last_seen":{"seconds":100}},
 {"id":2,"name":"phone","given_name":"phone","user":{"id":2,"name":"boris"},"ip_addresses":["100.64.0.2"],"online":true,"last_seen":{"seconds":200}},
 {"id":3,"name":"laptop","given_name":"laptop","user":{"id":2,"name":"boris"},"ip_addresses":["100.64.0.3"],"online":true},
 {"id":4,"name":"old","given_name":"old","user":{"id":2,"name":"boris"},"ip_addresses":["100.64.0.4"]}
]`

const statusJson = `{"BackendState":"Running","Peer":{
 "a":{"TailscaleIPs":["100.64.0.2"],"CurAddr":"9.9.9.9:41641","Relay":"syncloud","Online":true,"OS":"android"},
 "b":{"TailscaleIPs":["100.64.0.3"],"CurAddr":"","Relay":"syncloud","Online":true},
 "c":{"TailscaleIPs":["100.64.0.4"],"CurAddr":"","Relay":"syncloud","Online":false}}}`

func fakeDevices() (*Devices, *FakeExecutor) {
	executor := NewFakeExecutor()
	executor.outputs["hs nodes list -o json"] = nodesJson
	executor.outputs["ts status --json"] = statusJson
	executor.outputs["hs nodes delete"] = ""
	return NewDevices(executor, "hs", "ts", "device"), executor
}

func TestDevicesListShowsConnectionPath(t *testing.T) {
	devices, _ := fakeDevices()
	list, err := devices.List()
	require.NoError(t, err)
	require.Len(t, list, 4)

	assert.Equal(t, "device", list[0].Name)
	assert.True(t, list[0].Self)
	assert.Equal(t, "", list[0].Connection)

	byName := map[string]Device{}
	for _, device := range list {
		byName[device.Name] = device
	}
	assert.Equal(t, "direct", byName["phone"].Connection)
	assert.Equal(t, "9.9.9.9:41641", byName["phone"].Endpoint)
	assert.Equal(t, "relay", byName["laptop"].Connection)
	assert.Equal(t, "", byName["old"].Connection)
	assert.Equal(t, int64(200), byName["phone"].LastSeen)
	assert.Equal(t, "android", byName["phone"].Os)
	assert.Equal(t, "linux", list[0].Os)
}

func TestDevicesRemove(t *testing.T) {
	devices, executor := fakeDevices()
	require.NoError(t, devices.Remove(2))
	assert.Contains(t, executor.calls, "hs nodes delete --identifier 2 --force")
}

func TestDevicesRefusesToRemoveItself(t *testing.T) {
	devices, _ := fakeDevices()
	assert.ErrorIs(t, devices.Remove(1), ErrSelfDevice)
}

func TestDevicesRemoveUnknown(t *testing.T) {
	devices, _ := fakeDevices()
	assert.ErrorIs(t, devices.Remove(42), ErrUnknownDevice)
}
