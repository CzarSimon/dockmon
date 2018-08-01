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
	r := NewRouter(env.config.username, env.config.password)
	r.ServeDir("/", "frontend/web/build")
	r.GET("/health", handleHealthCheck, false)
	r.POST("/api/login", handleHealthCheck, true)
	r.GET("/api/status", env.getServiceStatus, true)
	r.GET("/api/statuses", env.getServiceStatuses, true)

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
