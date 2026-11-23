package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var ErrSelfDevice = errors.New("this device cannot be removed")
var ErrUnknownDevice = errors.New("unknown device")

type Devices struct {
	executor  Executor
	headscale string
	tailscale string
	hostname  string
}

func NewDevices(executor Executor, headscale, tailscale, hostname string) *Devices {
	return &Devices{executor: executor, headscale: headscale, tailscale: tailscale, hostname: hostname}
}

func (d *Devices) List() ([]Device, error) {
	nodes, err := d.nodes()
	if err != nil {
		return nil, err
	}
	peers := d.peers()
	devices := []Device{}
	for _, node := range nodes {
		device := Device{
			Id:        node.Id,
			Name:      node.GivenName,
			User:      node.User.Name,
			Addresses: node.IpAddresses,
			Online:    node.Online,
			LastSeen:  node.LastSeen.Seconds,
			Self:      d.isSelf(node),
		}
		if device.Name == "" {
			device.Name = node.Name
		}
		if !device.Self {
			device.Connection, device.Endpoint = connection(peers, node)
		}
		devices = append(devices, device)
	}
	sort.SliceStable(devices, func(i, j int) bool {
		if devices[i].Self != devices[j].Self {
			return devices[i].Self
		}
		return devices[i].Name < devices[j].Name
	})
	return devices, nil
}

func (d *Devices) Remove(id uint64) error {
	nodes, err := d.nodes()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if node.Id != id {
			continue
		}
		if d.isSelf(node) {
			return ErrSelfDevice
		}
		_, err := d.executor.Run(d.headscale, "nodes", "delete", "--identifier", fmt.Sprint(id), "--force")
		return err
	}
	return ErrUnknownDevice
}

func (d *Devices) isSelf(node HeadscaleNode) bool {
	return node.User.Name == DeviceUser && node.Name == d.hostname
}

func (d *Devices) nodes() ([]HeadscaleNode, error) {
	output, err := d.executor.Run(d.headscale, "nodes", "list", "-o", "json")
	if err != nil {
		return nil, err
	}
	var nodes []HeadscaleNode
	if err := json.Unmarshal(output, &nodes); err != nil {
		return nil, fmt.Errorf("headscale nodes: %w: %s", err, output)
	}
	return nodes, nil
}

func (d *Devices) peers() []PeerStatus {
	output, err := d.executor.Run(d.tailscale, "status", "--json")
	if err != nil {
		return nil
	}
	status := &Status{}
	if err := json.Unmarshal(output, status); err != nil {
		return nil
	}
	peers := []PeerStatus{}
	for _, peer := range status.Peer {
		peers = append(peers, peer)
	}
	return peers
}

func connection(peers []PeerStatus, node HeadscaleNode) (string, string) {
	for _, peer := range peers {
		if !sharesAddress(peer.TailscaleIPs, node.IpAddresses) {
			continue
		}
		if peer.CurAddr != "" {
			return "direct", peer.CurAddr
		}
		if peer.Relay != "" && node.Online {
			return "relay", peer.Relay
		}
		return "", ""
	}
	return "", ""
}

func sharesAddress(left, right []string) bool {
	for _, l := range left {
		for _, r := range right {
			if l == r {
				return true
			}
		}
	}
	return false
}
