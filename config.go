package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	endpoint "github.com/CzarSimon/go-endpoint"
	yaml "gopkg.in/yaml.v2"
)

const configFilename = "serviceConf.yml"

const (
	SERVICE_PORT = "DOCKMON_PORT"
	DB_NAME      = "DOCKMON_DB"
	DefaultPort  = "7777"
)

// LivenessOptions configuration options for a LivenessTarget.
type LivenessOptions struct {
	ServiceName      string `yaml:"serviceName" json:"serviceName"`
	LivenessURL      string `yaml:"livenessUrl" json:"livenessUrl"`
	LivenessInterval int    `yaml:"livenessInterval" json:"livenessInterval"`
	Restart          bool   `yaml:"restart" json:"restart"`
	FailAfter        uint8  `yaml:"failAfter" json:"failAfter"`
}

// config holds configuration options.
type config struct {
	serviceOptions []LivenessOptions
	port           string
	db             endpoint.SQLConfig
	httpTimeout    time.Duration
	dockerTimeout  time.Duration
}

// getConfig gets configuraton from both the environent and the serviceConf file.
func getConfig() config {
	serviceOptions, err := readServiceOptions(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	return config{
		serviceOptions: serviceOptions,
		port:           getServicePort(),
		db:             endpoint.NewPGConfig(DB_NAME),
		httpTimeout:    1 * time.Second,
		dockerTimeout:  10 * time.Second,
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
func readServiceOptions(filename string) ([]LivenessOptions, error) {
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var serviceOptions []LivenessOptions
	err = yaml.Unmarshal(rawData, &serviceOptions)
	return serviceOptions, err
}
