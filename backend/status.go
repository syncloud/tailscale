package main

type Status struct {
	BackendState string                `json:"BackendState"`
	Peer         map[string]PeerStatus `json:"Peer"`
}
