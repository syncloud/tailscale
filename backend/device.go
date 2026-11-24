package main

type Device struct {
	Id         uint64   `json:"id"`
	Name       string   `json:"name"`
	User       string   `json:"user"`
	Addresses  []string `json:"addresses"`
	Online     bool     `json:"online"`
	LastSeen   int64    `json:"last_seen"`
	Self       bool     `json:"self"`
	Connection string   `json:"connection"`
	Endpoint   string   `json:"endpoint"`
	Os         string   `json:"os"`
}
