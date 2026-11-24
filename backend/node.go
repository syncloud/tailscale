package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const DeviceUser = "device"

type Node struct {
	executor    Executor
	headscale   string
	tailscale   string
	loginServer string
	hostname    string
}

func NewNode(executor Executor, headscale, tailscale, loginServer, hostname string) *Node {
	return &Node{
		executor:    executor,
		headscale:   headscale,
		tailscale:   tailscale,
		loginServer: loginServer,
		hostname:    hostname,
	}
}

func (n *Node) Ensure() error {
	running, err := n.running()
	if err != nil {
		return err
	}
	if running {
		return nil
	}
	userId, err := n.deviceUser()
	if err != nil {
		return err
	}
	output, err := n.executor.Run(n.headscale, "preauthkeys", "create", "--user", fmt.Sprint(userId), "--expiration", "10m", "-o", "json")
	if err != nil {
		return err
	}
	key := &PreAuthKey{}
	if err := json.Unmarshal(output, key); err != nil {
		return fmt.Errorf("preauth key: %w: %s", err, output)
	}
	_, err = n.executor.Run(n.tailscale, "up", "--reset",
		"--login-server="+n.loginServer,
		"--authkey="+key.Key,
		"--hostname="+n.hostname,
		"--accept-dns=false",
		"--netfilter-mode=off",
		"--timeout=60s",
	)
	return err
}

func (n *Node) Ips() ([]string, error) {
	output, err := n.executor.Run(n.tailscale, "ip")
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(output)), nil
}

func (n *Node) running() (bool, error) {
	output, err := n.executor.Run(n.tailscale, "status", "--json")
	if err != nil {
		return false, nil
	}
	status := &Status{}
	if err := json.Unmarshal(output, status); err != nil {
		return false, fmt.Errorf("tailscale status: %w", err)
	}
	return status.BackendState == "Running", nil
}

func (n *Node) deviceUser() (uint64, error) {
	output, err := n.executor.Run(n.headscale, "users", "list", "-o", "json")
	if err != nil {
		return 0, err
	}
	var users []HeadscaleUser
	if err := json.Unmarshal(output, &users); err != nil {
		return 0, fmt.Errorf("headscale users: %w: %s", err, output)
	}
	for _, user := range users {
		if user.Name == DeviceUser {
			return user.Id, nil
		}
	}
	output, err = n.executor.Run(n.headscale, "users", "create", DeviceUser, "-o", "json")
	if err != nil {
		return 0, err
	}
	created := &HeadscaleUser{}
	if err := json.Unmarshal(output, created); err != nil {
		return 0, fmt.Errorf("headscale user create: %w: %s", err, output)
	}
	return created.Id, nil
}
