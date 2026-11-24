package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionRoundTrip(t *testing.T) {
	codec := NewSessionCodec([]byte("key"))
	value, err := codec.Encode(Session{User: "boris", Expires: time.Now().Add(time.Hour).Unix()})
	require.NoError(t, err)
	session, err := codec.Decode(value)
	require.NoError(t, err)
	assert.Equal(t, "boris", session.User)
}

func TestSessionRejectsOtherKey(t *testing.T) {
	value, err := NewSessionCodec([]byte("key")).Encode(Session{User: "boris", Expires: time.Now().Add(time.Hour).Unix()})
	require.NoError(t, err)
	_, err = NewSessionCodec([]byte("other")).Decode(value)
	assert.Error(t, err)
}

func TestSessionRejectsExpired(t *testing.T) {
	codec := NewSessionCodec([]byte("key"))
	value, err := codec.Encode(Session{User: "boris", Expires: time.Now().Add(-time.Minute).Unix()})
	require.NoError(t, err)
	_, err = codec.Decode(value)
	assert.Error(t, err)
}
