package main

type PeerStatus struct {
	TailscaleIPs []string `json:"TailscaleIPs"`
	CurAddr      string   `json:"CurAddr"`
	Relay        string   `json:"Relay"`
	Online       bool     `json:"Online"`
	OS           string   `json:"OS"`
}
