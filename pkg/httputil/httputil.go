package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// sendJSON Marshals a json body and sends as response.
func SendJSON(w http.ResponseWriter, v interface{}) (error, int) {
	js, err := json.Marshal(v)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return nil, http.StatusOK
}

// Router wrapper around a http.ServeMux to provide
// authentication for specific routes, mathing a routes to http methods
// and wrapping HandlerFuncs with error handling an logging.
type Router struct {
	mux      *http.ServeMux
	username string
	password string
}

// NewRouter creats a new Router with the given authentication credentials.
func NewRouter(username, password string) *Router {
	return &Router{
		mux:      http.NewServeMux(),
		username: username,
		password: password,
	}
}

// ServeHTTP passes each request to the underlying ServeMux.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

// GET wraps a HandlerFunc into a handler with optional authentication and
// registers it agains a GET method and pattern.
func (router *Router) GET(pattern string, h HandlerFunc, useAuth bool) {
	handler := NewHandler(http.MethodGet, router.username, router.password, h, useAuth)
	router.mux.Handle(pattern, handler)
}

// POST wraps a HandlerFunc into a handler with optional authentication and
// registers it agains a POST method and pattern.
func (router *Router) POST(pattern string, h HandlerFunc, useAuth bool) {
	handler := NewHandler(http.MethodPost, router.username, router.password, h, useAuth)
	router.mux.Handle(pattern, handler)
}

// ServeDir registers serviing of static files from a given directory.
func (router *Router) ServeDir(pattern, directory string) {
	router.mux.Handle(pattern, http.FileServer(http.Dir(directory)))
}

// HandlerFunc signature of a request handler.
type HandlerFunc func(http.ResponseWriter, *http.Request) (error, int)

// Handler wrapper around a HandlerFunc to provide
// authentication, method checking, logging and error handling.
type Handler struct {
	allowedMethod string
	handle        HandlerFunc
	username      string
	password      string
	useAuth       bool
}

// NewHandler creates and returns a new Handler.
func NewHandler(method, username, password string, h HandlerFunc, useAuth bool) Handler {
	return Handler{
		allowedMethod: method,
		handle:        h,
		useAuth:       useAuth,
		username:      username,
		password:      password,
	}
}

// ServeHTTP wrapps the call to the handlers HandlerFunc.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != h.allowedMethod {
		http.Error(w, fmt.Sprintf("Method %s not allowed\n", r.Method), http.StatusMethodNotAllowed)
		return
	}

	err := h.Authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	startTime := time.Now()
	err, status := h.handle(w, r)
	LogRequest(r, status, startTime)

	if err != nil {
		log.Printf("%d ERROR: %s\n", status, err)
		http.Error(w, err.Error(), status)
	}
}

// LogRequest logs: status, called route, method and completion time of a request.
func LogRequest(r *http.Request, status int, startTime time.Time) {
	var MilliPerNano int64 = 1000000
	requestTimeMS := time.Since(startTime).Nanoseconds() / MilliPerNano
	log.Printf("| %d | %s - %s | %d ms\n", status, r.Method, r.RequestURI, requestTimeMS)
}

// ParseQuery attempts to extract a query from
func ParseQuery(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return value, fmt.Errorf("No value found for key: %s", key)
	}
	return value, nil
}

// Authenticate checks the request credentials if the handler is set to do so.
func (h Handler) Authenticate(r *http.Request) error {
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
