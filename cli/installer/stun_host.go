package installer

import "strings"

var redirectDomains = []string{"syncloud.it", "syncloud.info"}

func StunHost(deviceDomain string) string {
	for _, redirect := range redirectDomains {
		if strings.HasSuffix(deviceDomain, "."+redirect) {
			return "stun." + redirect
		}
	}
	return "stun.syncloud.it"
}
