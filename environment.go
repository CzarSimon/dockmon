package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	docker "docker.io/go-docker"
	"github.com/CzarSimon/dockmon/pkg/datastore"
	"github.com/CzarSimon/dockmon/pkg/schema"
	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"
)

// Env holds reverence to environment variables and objects.
type Env struct {
	sigChan      chan os.Signal
	httpClient   *http.Client
	dockerClient *docker.Client
	serviceRepo  datastore.ServiceRepository
	config
}

// SetupEnv sets up an environment based on the current config.
func SetupEnv(config config) *Env {
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

func newServiceRepository(config config) datastore.ServiceRepository {
	db, err := config.db.Connect()
	failOnError(err)
	err = migrateDB(config.dbDriver, db)
	failOnError(err)

	serviceRepo := datastore.GetServiceRepository(config.dbDriver, db)

	for _, serviceOption := range config.serviceOptions {
		err := serviceRepo.SaveService(schema.NewServiceStatus(serviceOption))
		failOnError(err)
	}
	return serviceRepo
}

func migrateDB(dbDriver string, db *sql.DB) error {
	fmt.Printf("DB Driver: %s\n", dbDriver)
	migrationSource := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./resources/database"),
		Dir: dbDriver,
	}
	migrate.SetTable("migrations")

	migrations, err := migrationSource.FindMigrations()
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		return errors.New("Missing database migrations")
	}

	_, err = migrate.Exec(db, dbDriver, migrationSource, migrate.Up)
	if err != nil {
		return fmt.Errorf("Error applying database migrations: %s", err)
	}
	return nil
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
