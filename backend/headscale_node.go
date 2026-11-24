package main

type HeadscaleNode struct {
	Id          uint64        `json:"id"`
	Name        string        `json:"name"`
	GivenName   string        `json:"given_name"`
	User        HeadscaleUser `json:"user"`
	IpAddresses []string      `json:"ip_addresses"`
	Online      bool          `json:"online"`
	LastSeen    Timestamp     `json:"last_seen"`
}
