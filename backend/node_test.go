package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeRunningDoesNothing(t *testing.T) {
	executor := NewFakeExecutor()
	executor.outputs["ts status --json"] = `{"BackendState":"Running"}`
	require.NoError(t, NewNode(executor, "hs", "ts", "http://127.0.0.1:18731", "device").Ensure())
	assert.Equal(t, []string{"ts status --json"}, executor.calls)
}

func TestNodeJoinsWithNewKeyForExistingUser(t *testing.T) {
	executor := NewFakeExecutor()
	executor.outputs["ts status --json"] = `{"BackendState":"NeedsLogin"}`
	executor.outputs["hs users list -o json"] = `[{"id":3,"name":"alice"},{"id":7,"name":"device"}]`
	executor.outputs["hs preauthkeys create --user 7"] = `{"key":"hskey-auth-abc"}`
	executor.outputs["ts up"] = ""
	require.NoError(t, NewNode(executor, "hs", "ts", "http://127.0.0.1:18731", "device").Ensure())
	assert.Contains(t, executor.calls, "ts up --reset --login-server=http://127.0.0.1:18731 --authkey=hskey-auth-abc --hostname=device --accept-dns=false --netfilter-mode=off --timeout=60s")
}

func TestNodeCreatesDeviceUser(t *testing.T) {
	executor := NewFakeExecutor()
	executor.outputs["ts status --json"] = `{"BackendState":"NeedsLogin"}`
	executor.outputs["hs users list -o json"] = `[]`
	executor.outputs["hs users create device -o json"] = `{"id":1,"name":"device"}`
	executor.outputs["hs preauthkeys create --user 1"] = `{"key":"k"}`
	executor.outputs["ts up"] = ""
	require.NoError(t, NewNode(executor, "hs", "ts", "http://127.0.0.1:18731", "device").Ensure())
	assert.Contains(t, executor.calls, "hs users create device -o json")
}

func TestNodeJoinsWhenTailscaledNotAnswering(t *testing.T) {
	executor := NewFakeExecutor()
	executor.outputs["ts status --json"] = ""
	executor.failing["ts status --json"] = true
	executor.outputs["hs users list -o json"] = `[{"id":1,"name":"device"}]`
	executor.outputs["hs preauthkeys create --user 1"] = `{"key":"k"}`
	executor.outputs["ts up"] = ""
	require.NoError(t, NewNode(executor, "hs", "ts", "http://127.0.0.1:18731", "device").Ensure())
}

func TestNodeIps(t *testing.T) {
	executor := NewFakeExecutor()
	executor.outputs["ts ip"] = "100.64.0.1\nfd7a:115c:a1e0::1\n"
	ips, err := NewNode(executor, "hs", "ts", "", "").Ips()
	require.NoError(t, err)
	assert.Equal(t, []string{"100.64.0.1", "fd7a:115c:a1e0::1"}, ips)
}
