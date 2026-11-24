package main

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"sort"
)

var notApps = map[string]bool{
	"bare":      true,
	"bin":       true,
	"core":      true,
	"core18":    true,
	"core20":    true,
	"core22":    true,
	"core24":    true,
	"platform":  true,
	"snapd":     true,
	"tailscale": true,
}

type Records struct {
	file    string
	domain  string
	snapDir string
}

func NewRecords(file, domain, snapDir string) *Records {
	return &Records{file: file, domain: domain, snapDir: snapDir}
}

func (r *Records) Sync(ips []string) error {
	names, err := r.names()
	if err != nil {
		return err
	}
	records := []Record{}
	for _, name := range names {
		for _, ip := range ips {
			parsed := net.ParseIP(ip)
			if parsed == nil {
				continue
			}
			recordType := "AAAA"
			if parsed.To4() != nil {
				recordType = "A"
			}
			records = append(records, Record{Name: name, Type: recordType, Value: ip})
		}
	}
	content, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	existing, err := os.ReadFile(r.file)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}
	tmp := r.file + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.file)
}

func (r *Records) names() ([]string, error) {
	entries, err := os.ReadDir(r.snapDir)
	if err != nil {
		return nil, err
	}
	names := []string{r.domain, "auth." + r.domain}
	for _, entry := range entries {
		if !entry.IsDir() || notApps[entry.Name()] {
			continue
		}
		names = append(names, entry.Name()+"."+r.domain)
	}
	sort.Strings(names[2:])
	return names, nil
}
