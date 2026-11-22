package main

import (
	"flag"
	"log"
	"net/url"
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

	executor := &ShellExecutor{}
	node := NewNode(executor, config.Headscale, config.Tailscale, config.HeadscaleUrl, config.Hostname)
	records := NewRecords(config.Records, config.Domain, config.SnapDir)
	go func() {
		for {
			sync(node, records)
			time.Sleep(30 * time.Second)
		}
	}()

	log.Printf("listening on %s", config.Socket)
	log.Fatal(NewProxy(config.Socket, upstream, NewPage(config.Url)).Start())
}

func sync(node *Node, records *Records) {
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
