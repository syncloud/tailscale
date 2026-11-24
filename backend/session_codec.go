package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type SessionCodec struct {
	key []byte
	now func() time.Time
}

func NewSessionCodec(key []byte) *SessionCodec {
	return &SessionCodec{key: key, now: time.Now}
}

func (c *SessionCodec) Encode(session Session) (string, error) {
	body, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(body) + "." + base64.RawURLEncoding.EncodeToString(c.sign(body)), nil
}

func (c *SessionCodec) Decode(value string) (*Session, error) {
	parts := strings.SplitN(value, ".", 2)
	if len(parts) != 2 {
		return nil, errors.New("malformed session")
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(signature, c.sign(body)) {
		return nil, errors.New("bad session signature")
	}
	session := &Session{}
	if err := json.Unmarshal(body, session); err != nil {
		return nil, err
	}
	if c.now().Unix() >= session.Expires {
		return nil, errors.New("session expired")
	}
	return session, nil
}

func (c *SessionCodec) sign(body []byte) []byte {
	mac := hmac.New(sha256.New, c.key)
	mac.Write(body)
	return mac.Sum(nil)
}
