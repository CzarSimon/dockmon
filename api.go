package main

import (
	"log"
	"net/http"
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
	r := http.NewServeMux()
	wrap := newWrapper(env.config.username, env.config.password)

	r.Handle("/", http.FileServer(http.Dir("frontend/web/build")))
	r.Handle("/health", wrap.new(http.MethodGet, handleHealthCheck))
	r.Handle("/api/login", wrap.newWithAuth(http.MethodPost, handleHealthCheck))
	r.Handle("/api/status", wrap.newWithAuth(http.MethodGet, env.getServiceStatus))
	r.Handle("/api/statuses", wrap.newWithAuth(http.MethodGet, env.getServiceStatuses))

	return &http.Server{
		Addr:    ":" + env.port,
		Handler: r,
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
