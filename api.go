package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// startAPI start the monitoring agents rest api.
func (env *Env) startAPI() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", handleHealthCheck)
	err := e.Start(":" + env.port)
	if err != nil {
		log.Println(err)
	}
}

// handleHealthCheck returns a 200 OK on being invoked.
func handleHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
