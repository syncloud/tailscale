package main

type Session struct {
	User    string `json:"user"`
	Expires int64  `json:"exp"`
}
