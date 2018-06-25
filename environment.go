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
	serviceRepo  ServiceRepository
	config
}

// SetupEnv sets up an environment based on the current config.
func SetupEnv(config config) *Env {
	log.Println(config.db)
	return &Env{
		sigChan:      make(chan os.Signal),
		httpClient:   newHttpClient(config),
		dockerClient: newDockerClient(),
		serviceRepo:  newServiceRepository(config),
		config:       config,
	}
}

// Close close relevant pointers in the environment.
func (env *Env) Close() error {
	return env.serviceRepo.Close()
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
	failOnError(err)
	return client
}

func newServiceRepository(config config) ServiceRepository {
	db, err := config.db.Connect()
	failOnError(err)
	serviceRepo := &PgServiceRepo{db: db}

	for _, serviceOption := range config.serviceOptions {
		err := serviceRepo.SaveService(NewServiceStatus(serviceOption))
		failOnError(err)
	}
	return serviceRepo
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
