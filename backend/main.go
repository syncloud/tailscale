package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	configFile := flag.String("config", "", "config file")
	flag.Parse()

	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	upstream, err := url.Parse(config.HeadscaleUrl)
	if err != nil {
		log.Fatal(err)
	}
	authUrl, err := url.Parse(config.AuthUrl)
	if err != nil {
		log.Fatal(err)
	}
	sessionKey, err := os.ReadFile(config.SessionKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	executor := &ShellExecutor{}
	node := NewNode(executor, config.Headscale, config.Tailscale, config.HeadscaleUrl, config.Hostname)
	records := NewRecords(config.Records, config.Domain, config.SnapDir)
	go func() {
		for {
			syncDevice(node, records)
			time.Sleep(30 * time.Second)
		}
	}()

	client := &http.Client{Transport: NewAuthTransport(authUrl.Host, config.AuthSocket), Timeout: 30 * time.Second}
	secret := func() (string, error) {
		content, err := os.ReadFile(config.SecretFile)
		return strings.TrimSpace(string(content)), err
	}
	auth := NewOidc(config.AuthUrl, config.ClientId, secret, config.Url+CallbackPath, client, NewSessionCodec(sessionKey))
	auth.InitWithRetry(context.Background(), 30*time.Second)

	devices := NewDevices(executor, config.Headscale, config.Tailscale, config.Hostname)
	router := NewRouter(config.Url, config.WebDir, auth, devices, NewHeadscaleProxy(upstream))
	log.Printf("listening on %s", config.Socket)
	log.Fatal(router.Start(config.Socket))
}

func syncDevice(node *Node, records *Records) {
	if err := node.Ensure(); err != nil {
		log.Printf("device node: %v", err)
		return
	}
	ips, err := node.Ips()
	if err != nil {
		log.Printf("device node ips: %v", err)
		return
	}
	if err := records.Sync(ips); err != nil {
		log.Printf("dns records: %v", err)
	}
}
