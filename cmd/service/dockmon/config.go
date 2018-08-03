package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
	endpoint "github.com/CzarSimon/go-endpoint"
	yaml "gopkg.in/yaml.v2"
)

const configFilename = "serviceConf.yml"

const (
	SERVICE_PORT       = "DOCKMON_PORT"
	DB_NAME            = "DOCKMON_DB"
	USERNAME_KEY       = "DOCKMON_USERNAME"
	PASSWORD_KEY       = "DOCKMON_PASSWORD"
	DefaultPort        = "7777"
	STORAGE_FLAG       = "storage"
	DefaultStorageType = "postgres"
)

// config holds configuration options.
type config struct {
	serviceOptions []schema.LivenessOptions
	port           string
	db             endpoint.SQLConfig
	dbDriver       string
	httpTimeout    time.Duration
	dockerTimeout  time.Duration
	username       string
	password       string
}

// getConfig gets configuraton from both the environent and the serviceConf file.
func getConfig() config {
	serviceOptions, err := readServiceOptions(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	dbConfig := getDBConfig()

	return config{
		serviceOptions: serviceOptions,
		port:           getServicePort(),
		db:             dbConfig,
		dbDriver:       dbConfig.ConnInfo().DriverName,
		httpTimeout:    1 * time.Second,
		dockerTimeout:  10 * time.Second,
		username:       os.Getenv(USERNAME_KEY),
		password:       os.Getenv(PASSWORD_KEY),
	}
}

func getDBConfig() endpoint.SQLConfig {
	var dbDriver string
	flag.StringVar(&dbDriver, STORAGE_FLAG, DefaultStorageType, "Storage type to use")
	flag.Parse()

	switch dbDriver {
	case "postgres":
		return endpoint.NewPGConfig(DB_NAME)
	case "mysql":
		return endpoint.NewMySQLConfig(DB_NAME)
	case "sqlite3":
		return endpoint.NewSQLiteConfig(DB_NAME)
	case "memory":
		return endpoint.NewSQLiteConfig(":memory:")
	default:
		return endpoint.NewPGConfig(DB_NAME)
	}
}

// getServicePort gets the port on which to serve the rest api.
func getServicePort() string {
	port := os.Getenv(SERVICE_PORT)
	if port == "" {
		return DefaultPort
	}
	return port
}

// readServiceOptions reads service options in the provided serviceConf.yml file.
func readServiceOptions(filename string) ([]schema.LivenessOptions, error) {
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var serviceOptions []schema.LivenessOptions
	err = yaml.Unmarshal(rawData, &serviceOptions)
	return serviceOptions, err
}
