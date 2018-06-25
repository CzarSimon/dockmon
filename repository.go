package main

import (
	"database/sql"
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
	Close() error
}

// PgServiceRepo postgres implementation of the ServiceRepository interface.
type PgServiceRepo struct {
	db *sql.DB
}

const insertServiceStatusQuery = `
  INSERT INTO dockmon_liveness_target (
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT DO NOTHING`

// SaveService inserts a new ServiceStatus into the database.
func (repo *PgServiceRepo) SaveService(serviceStatus ServiceStatus) error {
	stmt, err := repo.db.Prepare(insertServiceStatusQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		serviceStatus.ServiceName, serviceStatus.LivenessURL, serviceStatus.LivenessInterval,
		serviceStatus.ShouldRestart, serviceStatus.FailAfter, serviceStatus.IsHealty,
		serviceStatus.Restarts, serviceStatus.ConsecutiveFailedHealthChecks,
		serviceStatus.LastRestarted, serviceStatus.LastHealthSuccess,
		serviceStatus.LastHealthFailure, serviceStatus.CreatedAt)
	return err
}

const selectServiceStatusQuery = `
  SELECT
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at
  FROM dockmon_liveness_target WHERE service_name = $1`

// GetServiceStatus gets a specified service status from the database.
func (repo *PgServiceRepo) GetServiceStatus(serviceName string) (ServiceStatus, error) {
	var s ServiceStatus
	err := repo.db.QueryRow(selectServiceStatusQuery, serviceName).Scan(
		&s.ServiceName, &s.LivenessURL, &s.LivenessInterval, &s.ShouldRestart,
		&s.FailAfter, &s.IsHealty, &s.Restarts, &s.ConsecutiveFailedHealthChecks,
		&s.LastRestarted, &s.LastHealthSuccess, &s.LastHealthFailure, &s.CreatedAt)
	if err != nil {
		return emptyServiceStatus, err
	}
	return s, err
}

const selectServiceStatusesQuery = `
  SELECT
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at
  FROM dockmon_liveness_target`

// GetServiceStatuses gets all service statuses from the database.
func (repo *PgServiceRepo) GetServiceStatuses() ([]ServiceStatus, error) {
	rows, err := repo.db.Query(selectServiceStatusesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return createServiceStatusesFromRows(rows)
}

// createServiceStatusesFromRows turns a resulting list of rows into
// a list of service statuses.
func createServiceStatusesFromRows(rows *sql.Rows) ([]ServiceStatus, error) {
	statuses := make([]ServiceStatus, 0)
	var s ServiceStatus
	for rows.Next() {
		err := rows.Scan(
			&s.ServiceName, &s.LivenessURL, &s.LivenessInterval, &s.ShouldRestart,
			&s.FailAfter, &s.IsHealty, &s.Restarts, &s.ConsecutiveFailedHealthChecks,
			&s.LastRestarted, &s.LastHealthSuccess, &s.LastHealthFailure, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

// Close closes the underlying database connection.
func (repo *PgServiceRepo) Close() error {
	return repo.Close()
}
