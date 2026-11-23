package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(url, webDir string, auth *Oidc, devices *Devices, headscale http.Handler) *Router {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, path.Join(webDir, "index.html"))
	})
	mux.Handle("GET /app/assets/", http.StripPrefix("/app", http.FileServer(http.Dir(webDir))))
	mux.HandleFunc("GET "+LoginPath, auth.Login)
	mux.HandleFunc("GET "+CallbackPath, auth.Callback)
	mux.HandleFunc("GET "+LogoutPath, auth.Logout)
	mux.HandleFunc("GET /app/api/me", secured(auth, func(w http.ResponseWriter, r *http.Request, session *Session) {
		writeJson(w, http.StatusOK, Me{User: session.User, Url: url})
	}))
	mux.HandleFunc("GET /app/api/devices", secured(auth, func(w http.ResponseWriter, r *http.Request, _ *Session) {
		list, err := devices.List()
		if err != nil {
			log.Printf("devices: %v", err)
			writeJson(w, http.StatusInternalServerError, ApiError{Error: "unable to list devices"})
			return
		}
		writeJson(w, http.StatusOK, list)
	}))
	mux.HandleFunc("DELETE /app/api/devices/{id}", secured(auth, func(w http.ResponseWriter, r *http.Request, session *Session) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			writeJson(w, http.StatusBadRequest, ApiError{Error: "bad device id"})
			return
		}
		err = devices.Remove(id)
		switch {
		case errors.Is(err, ErrSelfDevice):
			writeJson(w, http.StatusForbidden, ApiError{Error: err.Error()})
		case errors.Is(err, ErrUnknownDevice):
			writeJson(w, http.StatusNotFound, ApiError{Error: err.Error()})
		case err != nil:
			log.Printf("remove device %d: %v", id, err)
			writeJson(w, http.StatusInternalServerError, ApiError{Error: "unable to remove device"})
		default:
			log.Printf("device %d removed by %s", id, session.User)
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	mux.Handle("/", headscale)
	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Start(socket string) error {
	_ = os.Remove(socket)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}
	if err := os.Chmod(socket, 0666); err != nil {
		return err
	}
	return http.Serve(listener, r)
}

func secured(auth *Oidc, handler func(http.ResponseWriter, *http.Request, *Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := auth.Session(r)
		if !ok {
			writeJson(w, http.StatusUnauthorized, ApiError{Error: "sign in required"})
			return
		}
		handler(w, r, session)
	}
}

func writeJson(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
