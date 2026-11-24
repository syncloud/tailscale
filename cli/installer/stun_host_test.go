package installer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStunHostFollowsRedirectDomain(t *testing.T) {
	assert.Equal(t, "stun.syncloud.it", StunHost("borisarm64.syncloud.it"))
	assert.Equal(t, "stun.syncloud.info", StunHost("test.syncloud.info"))
}

func TestStunHostDefaultsToProdForCustomDomains(t *testing.T) {
	assert.Equal(t, "stun.syncloud.it", StunHost("bookworm.com"))
	assert.Equal(t, "stun.syncloud.it", StunHost("notsyncloud.it"))
}
