package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// startAPI start the monitoring agents rest api.
func (env *Env) startAPI() {
	server := registerRoutes(env)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// registerRoutes registers api routes.
func registerRoutes(env *Env) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("frontend/web/build")))
	mux.Handle("/health", wrapHandle(http.MethodGet, handleHealthCheck))
	mux.Handle("/api/status", wrapHandle(http.MethodGet, env.getServiceStatus))
	mux.Handle("/api/statuses", wrapHandle(http.MethodGet, env.getServiceStatuses))

	return &http.Server{
		Addr:    ":" + env.port,
		Handler: mux,
	}
}

// getServiceStatus gets the health status of specified monitored service.
func (env *Env) getServiceStatus(w http.ResponseWriter, r *http.Request) (error, int) {
	serviceName, err := parseQuery(r, "serviceName")
	if err != nil {
		return err, http.StatusBadRequest
	}

	serviceStatus, err := env.serviceRepo.GetServiceStatus(serviceName)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return sendJSON(w, serviceStatus)
}

// getServiceStatuses gets the health status of all monitored services.
func (env *Env) getServiceStatuses(w http.ResponseWriter, r *http.Request) (error, int) {
	serviceStatuses, err := env.serviceRepo.GetServiceStatuses()
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return sendJSON(w, serviceStatuses)
}

// handleHealthCheck returns a 200 OK on being invoked.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) (error, int) {
	return sendJSON(w, map[string]string{"status": "OK"})
}

func sendError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, gin.H{"error": errorMsg})
}

// SendJSON Marshals a json body and sends as response.
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

type handlerFunc func(http.ResponseWriter, *http.Request) (error, int)

type handler struct {
	allowedMethod string
	handle        handlerFunc
}

func wrapHandle(method string, h handlerFunc) handler {
	return handler{
		allowedMethod: method,
		handle:        h,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != h.allowedMethod {
		http.Error(w, fmt.Sprintf("Method %s not allowed\n", r.Method), http.StatusMethodNotAllowed)
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
