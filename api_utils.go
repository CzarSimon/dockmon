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

type wrapper struct {
	username string
	password string
}

func newWrapper(username, password string) wrapper {
	return wrapper{
		username: username,
		password: password,
	}
}

type handlerFunc func(http.ResponseWriter, *http.Request) (error, int)

type handler struct {
	allowedMethod string
	handle        handlerFunc
	username      string
	password      string
	useAuth       bool
}

func (w wrapper) new(method string, h handlerFunc) handler {
	return handler{
		allowedMethod: method,
		handle:        h,
		useAuth:       false,
	}
}

func (w wrapper) newWithAuth(method string, h handlerFunc) handler {
	return handler{
		allowedMethod: method,
		handle:        h,
		useAuth:       true,
		username:      w.username,
		password:      w.password,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h handler) authenticate(r *http.Request) error {
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
