package installer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthSocketPathFromProxyPass(t *testing.T) {
	assert.Equal(t, "/var/snap/platform/current/authelia.socket", AuthSocketPath("http://unix:/var/snap/platform/current/authelia.socket:"))
}

func TestAuthSocketPathKeepsPlainPath(t *testing.T) {
	assert.Equal(t, "/run/authelia.socket", AuthSocketPath("/run/authelia.socket"))
}
