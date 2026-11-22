package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Socket       string `json:"socket"`
	HeadscaleUrl string `json:"headscale_url"`
	Headscale    string `json:"headscale"`
	Tailscale    string `json:"tailscale"`
	Domain       string `json:"domain"`
	AppDomain    string `json:"app_domain"`
	Url          string `json:"url"`
	Hostname     string `json:"hostname"`
	Records      string `json:"records"`
	SnapDir      string `json:"snap_dir"`
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
