package installer

import "strings"

func AuthSocketPath(proxyPass string) string {
	return strings.TrimSuffix(strings.TrimPrefix(proxyPass, "http://unix:"), ":")
}
