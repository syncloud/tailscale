package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const (
	stateCookie    = "tailscale_oidc_state"
	verifierCookie = "tailscale_oidc_verifier"
	sessionCookie  = "tailscale_session"
	sessionTtl     = 12 * time.Hour
	LoginPath      = "/app/auth/login"
	CallbackPath   = "/app/auth/callback"
	LogoutPath     = "/app/auth/logout"
)

type Oidc struct {
	issuer      string
	clientId    string
	secret      func() (string, error)
	redirectUrl string
	logoutUrl   string
	client      *http.Client
	codec       *SessionCodec
	mutex       sync.RWMutex
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
}

func NewOidc(issuer, clientId string, secret func() (string, error), redirectUrl string, client *http.Client, codec *SessionCodec) *Oidc {
	return &Oidc{
		issuer:      issuer,
		clientId:    clientId,
		secret:      secret,
		redirectUrl: redirectUrl,
		logoutUrl:   issuer + "/logout",
		client:      client,
		codec:       codec,
	}
}

func (o *Oidc) InitWithRetry(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			provider, err := oidc.NewProvider(o.context(ctx), o.issuer)
			if err == nil {
				o.mutex.Lock()
				o.provider = provider
				o.verifier = provider.Verifier(&oidc.Config{ClientID: o.clientId})
				o.mutex.Unlock()
				log.Printf("oidc ready: %s", o.issuer)
				return
			}
			log.Printf("oidc discovery: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
			}
		}
	}()
}

func (o *Oidc) Login(w http.ResponseWriter, r *http.Request) {
	config, ok := o.config()
	if !ok {
		http.Error(w, "sign-in is starting, try again shortly", http.StatusServiceUnavailable)
		return
	}
	state, err := random(16)
	if err != nil {
		http.Error(w, "state", http.StatusInternalServerError)
		return
	}
	verifier, err := random(32)
	if err != nil {
		http.Error(w, "verifier", http.StatusInternalServerError)
		return
	}
	setCookie(w, stateCookie, state, 10*time.Minute)
	setCookie(w, verifierCookie, verifier, 10*time.Minute)
	challenge := sha256.Sum256([]byte(verifier))
	http.Redirect(w, r, config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", base64.RawURLEncoding.EncodeToString(challenge[:])),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	), http.StatusFound)
}

func (o *Oidc) Callback(w http.ResponseWriter, r *http.Request) {
	config, ok := o.config()
	if !ok {
		http.Error(w, "sign-in is starting, try again shortly", http.StatusServiceUnavailable)
		return
	}
	state, err := r.Cookie(stateCookie)
	if err != nil || r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}
	verifier, err := r.Cookie(verifierCookie)
	if err != nil {
		http.Error(w, "verifier missing", http.StatusBadRequest)
		return
	}
	ctx := o.context(r.Context())
	token, err := config.Exchange(ctx, r.URL.Query().Get("code"), oauth2.SetAuthURLParam("code_verifier", verifier.Value))
	if err != nil {
		log.Printf("oidc exchange: %v", err)
		http.Error(w, "sign-in failed", http.StatusBadGateway)
		return
	}
	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "id token missing", http.StatusBadGateway)
		return
	}
	o.mutex.RLock()
	verifierOidc := o.verifier
	o.mutex.RUnlock()
	idToken, err := verifierOidc.Verify(ctx, rawIdToken)
	if err != nil {
		log.Printf("oidc verify: %v", err)
		http.Error(w, "sign-in failed", http.StatusUnauthorized)
		return
	}
	claims := Claims{}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "claims", http.StatusBadGateway)
		return
	}
	value, err := o.codec.Encode(Session{User: claims.User(), Expires: time.Now().Add(sessionTtl).Unix()})
	if err != nil {
		http.Error(w, "session", http.StatusInternalServerError)
		return
	}
	clearCookie(w, stateCookie)
	clearCookie(w, verifierCookie)
	setCookie(w, sessionCookie, value, sessionTtl)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (o *Oidc) Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, sessionCookie)
	http.Redirect(w, r, o.logoutUrl, http.StatusFound)
}

func (o *Oidc) Session(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil, false
	}
	session, err := o.codec.Decode(cookie.Value)
	if err != nil {
		return nil, false
	}
	return session, true
}

func (o *Oidc) config() (*oauth2.Config, bool) {
	o.mutex.RLock()
	provider := o.provider
	o.mutex.RUnlock()
	if provider == nil {
		return nil, false
	}
	secret, err := o.secret()
	if err != nil {
		log.Printf("oidc secret: %v", err)
		return nil, false
	}
	return &oauth2.Config{
		ClientID:     o.clientId,
		ClientSecret: secret,
		RedirectURL:  o.redirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}, true
}

func (o *Oidc) context(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, o.client)
}

func setCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Path: "/", MaxAge: -1, HttpOnly: true, Secure: true})
}

func random(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
