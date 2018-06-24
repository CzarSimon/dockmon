package main

import (
	"log"
	"net/http"
	"os"

	docker "docker.io/go-docker"
)

// Env holds reverence to environment variables and objects.
type Env struct {
	sigChan      chan os.Signal
	httpClient   *http.Client
	dockerClient *docker.Client
	config
}

// SetupEnv sets up an environment based on the current config.
func SetupEnv(config config) *Env {
	return &Env{
		sigChan:      make(chan os.Signal),
		httpClient:   newHttpClient(config),
		dockerClient: newDockerClient(),
		config:       config,
	}
}

// newHttpClient sets up a new http client with a configured timeout.
func newHttpClient(config config) *http.Client {
	return &http.Client{
		Timeout: config.httpTimeout,
	}
}

// newDockerClient sets up a new docker client and panics on error.
func newDockerClient() *docker.Client {
	client, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	return client
}
