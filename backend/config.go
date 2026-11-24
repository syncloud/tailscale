package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Socket         string `json:"socket"`
	HeadscaleUrl   string `json:"headscale_url"`
	Headscale      string `json:"headscale"`
	Tailscale      string `json:"tailscale"`
	Domain         string `json:"domain"`
	AppDomain      string `json:"app_domain"`
	Url            string `json:"url"`
	Hostname       string `json:"hostname"`
	Records        string `json:"records"`
	SnapDir        string `json:"snap_dir"`
	WebDir         string `json:"web_dir"`
	AuthUrl        string `json:"auth_url"`
	AuthSocket     string `json:"auth_socket"`
	ClientId       string `json:"client_id"`
	SecretFile     string `json:"secret_file"`
	SessionKeyFile string `json:"session_key_file"`
}

func LoadConfig(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}
	return config, nil
}
