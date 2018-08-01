package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// sendJSON Marshals a json body and sends as response.
func sendJSON(w http.ResponseWriter, v interface{}) (error, int) {
	js, err := json.Marshal(v)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return nil, http.StatusOK
}

type Router struct {
	mux      *http.ServeMux
	username string
	password string
}

func NewRouter(username, password string) *Router {
	return &Router{
		mux:      http.NewServeMux(),
		username: username,
		password: password,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) GET(pattern string, h HandlerFunc, useAuth bool) {
	handler := NewHandler(http.MethodGet, router.username, router.password, h, useAuth)
	router.mux.Handle(pattern, handler)
}

func (router *Router) POST(pattern string, h HandlerFunc, useAuth bool) {
	handler := NewHandler(http.MethodPost, router.username, router.password, h, useAuth)
	router.mux.Handle(pattern, handler)
}

func (router *Router) ServeDir(pattern, directory string) {
	router.mux.Handle("/", http.FileServer(http.Dir(directory)))
}

type HandlerFunc func(http.ResponseWriter, *http.Request) (error, int)

type Handler struct {
	allowedMethod string
	handle        HandlerFunc
	username      string
	password      string
	useAuth       bool
}

func NewHandler(method, username, password string, h HandlerFunc, useAuth bool) Handler {
	return Handler{
		allowedMethod: method,
		handle:        h,
		useAuth:       useAuth,
		username:      username,
		password:      password,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != h.allowedMethod {
		http.Error(w, fmt.Sprintf("Method %s not allowed\n", r.Method), http.StatusMethodNotAllowed)
		return
	}

	err := h.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	startTime := time.Now()
	err, status := h.handle(w, r)
	logRequest(r, status, startTime)

	if err != nil {
		log.Printf("%d ERROR: %s\n", status, err)
		http.Error(w, err.Error(), status)
	}
}

func logRequest(r *http.Request, status int, startTime time.Time) {
	var MilliPerNano int64 = 1000000
	requestTimeMS := time.Since(startTime).Nanoseconds() / MilliPerNano
	log.Printf("| %d | %s - %s | %d ms\n", status, r.Method, r.RequestURI, requestTimeMS)
}

func parseQuery(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return value, fmt.Errorf("No value found for key: %s", key)
	}
	return value, nil
}

func (h Handler) authenticate(r *http.Request) error {
	if !h.useAuth {
		return nil
	}

	username, password, ok := r.BasicAuth()
	if !ok || username != h.username || password != h.password {
		log.Printf("%d - Authentication failed for user: %s\n", http.StatusUnauthorized, username)
		return errors.New("User could not be authenticated")
	}
	return nil
}
