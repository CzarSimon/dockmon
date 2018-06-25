package main

import (
	"log"
	"net/http"

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
	r := gin.Default()
	r.GET("/health", handleHealthCheck)
	r.GET("/api/status/:serviceName", env.getServiceStatus)
	r.GET("/api/statuses", env.getServiceStatuses)

	return &http.Server{
		Addr:    ":" + env.port,
		Handler: r,
	}
}

// getServiceStatus gets the health status of specified monitored service.
func (env *Env) getServiceStatus(c *gin.Context) {
	serviceName := c.Param("serviceName")
	serviceStatus, err := env.serviceRepo.GetServiceStatus(serviceName)
	if err != nil {
		log.Println(err)
		sendError(c, http.StatusInternalServerError, "Error in retrieving service status")
		return
	}
	c.JSON(http.StatusOK, serviceStatus)
}

// getServiceStatuses gets the health status of all monitored services.
func (env *Env) getServiceStatuses(c *gin.Context) {
	serviceStatuses, err := env.serviceRepo.GetServiceStatuses()
	if err != nil {
		log.Println(err)
		sendError(c, http.StatusInternalServerError, "Error in retrieving service statuses")
		return
	}
	c.JSON(http.StatusOK, serviceStatuses)
}

// handleHealthCheck returns a 200 OK on being invoked.
func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func sendError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, gin.H{"error": errorMsg})
}
