package main

import (
	"time"

	_ "github.com/lib/pq"
)

var (
	emptyServiceStatus ServiceStatus = ServiceStatus{}
	beginingOfTime                   = time.Time{}.UTC()
)

// ServiceStatus contains metadata about a service and its health status and history.
type ServiceStatus struct {
	ServiceName                   string    `json:"serviceName"`
	LivenessURL                   string    `json:"livenessUrl"`
	LivenessInterval              int       `json:"livenessInterval"`
	ShouldRestart                 bool      `json:"shouldRestart"`
	FailAfter                     int       `json:"failAfter"`
	IsHealty                      bool      `json:"isHealty"`
	Restarts                      int       `json:"restarts"`
	ConsecutiveFailedHealthChecks int       `json:"consecutiveFailedHealthChecks"`
	LastRestarted                 time.Time `json:"lastRestarted"`
	LastHealthSuccess             time.Time `json:"lastHealthSuccess"`
	LastHealthFailure             time.Time `json:"lastHealthFailure"`
	CreatedAt                     time.Time `json:"createdAt"`
}

func NewServiceStatus(opts LivenessOptions) ServiceStatus {
	return ServiceStatus{
		ServiceName:                   opts.ServiceName,
		LivenessURL:                   opts.LivenessURL,
		LivenessInterval:              opts.LivenessInterval,
		ShouldRestart:                 opts.Restart,
		FailAfter:                     int(opts.FailAfter),
		IsHealty:                      true,
		Restarts:                      0,
		ConsecutiveFailedHealthChecks: 0,
		LastRestarted:                 beginingOfTime,
		LastHealthSuccess:             beginingOfTime,
		LastHealthFailure:             beginingOfTime,
		CreatedAt:                     time.Now().UTC(),
	}
}

// ServiceRepository interface to persist, update and retrieve service statuses.
type ServiceRepository interface {
	SaveService(serviceStatus ServiceStatus) error
	GetServiceStatus(serviceName string) (ServiceStatus, error)
	GetServiceStatuses() ([]ServiceStatus, error)

	SaveHealthSuccess(serviceName string, timestamp time.Time) error
	SaveHealthFailure(serviceName string, timestamp time.Time) error
	SaveRestart(serviceName string, timestamp time.Time) error
	Close() error
}
