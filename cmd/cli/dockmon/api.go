package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
)

// ApiClient interface for clients interacting with dockmon.
type ApiClient interface {
	GetStatuses() []schema.ServiceStatus
	GetStatus(serviceName string) schema.ServiceStatus
	Login()
}

// RESTApiClient client used to interact with the Dockmon REST api.
type RESTApiClient struct {
	config     Config
	httpClient *http.Client
}

// NewApiClient creates a new RESTApiClient.
func NewApiClient(config Config) ApiClient {
	return RESTApiClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func GetApiClientAndTestCredentials() ApiClient {
	config, err := getConfig()
	failOnError(err)
	api := NewApiClient(config)
	api.Login()
	return api
}

// GetStatuses gets the configured list of services along with each service status.
func (api RESTApiClient) GetStatuses() []schema.ServiceStatus {
	resp := api.performRequest(api.createGetRequest("/api/statuses"))
	defer resp.Body.Close()

	serviceStatuses := make([]schema.ServiceStatus, 0)
	err := json.NewDecoder(resp.Body).Decode(&serviceStatuses)
	failOnError(err)

	return serviceStatuses
}

// GetStatuses gets the a specific services along with its service status.
func (api RESTApiClient) GetStatus(serviceName string) schema.ServiceStatus {
	route := fmt.Sprintf("/api/status?serviceName=%s", serviceName)
	resp := api.performRequest(api.createGetRequest(route))
	defer resp.Body.Close()

	var serviceStatus schema.ServiceStatus
	err := json.NewDecoder(resp.Body).Decode(&serviceStatus)
	failOnError(err)

	return serviceStatus
}

// GetStatuses gets the a specific services along with its service status.
func (api RESTApiClient) Login() {
	resp := api.performRequest(api.createPostRequest("/api/login", nil))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Invalid login credentials")
		os.Exit(1)
	}
}

func (api RESTApiClient) performRequest(r *http.Request) *http.Response {
	resp, err := api.httpClient.Do(r)
	failOnError(err)

	return resp
}

func (api RESTApiClient) createGetRequest(route string) *http.Request {
	return api.createRequest(http.MethodGet, route, nil)
}

func (api RESTApiClient) createPostRequest(route string, body io.Reader) *http.Request {
	return api.createRequest(http.MethodPost, route, body)
}

func (api RESTApiClient) createRequest(method, route string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, api.makeURL(route), body)
	failOnError(err)
	api.setBasicAuth(req)
	return req
}

func (api RESTApiClient) setBasicAuth(r *http.Request) {
	credentials := fmt.Sprintf("%s:%s", api.config.Username, api.config.Password)
	encodedCreds := base64.StdEncoding.EncodeToString([]byte(credentials))

	r.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedCreds))
}

func (api RESTApiClient) makeURL(route string) string {
	return api.config.Host + route
}

func failOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
