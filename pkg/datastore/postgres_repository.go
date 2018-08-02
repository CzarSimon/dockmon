package datastore

import (
	"database/sql"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
	_ "github.com/lib/pq"
)

// PgServiceRepo postgres implementation of the ServiceRepository interface.
type PgServiceRepo struct {
	db *sql.DB
}

const pgInsertServiceStatusQuery = `
  INSERT INTO dockmon_liveness_target (
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT DO NOTHING`

// SaveService inserts a new ServiceStatus into the database.
func (repo *PgServiceRepo) SaveService(serviceStatus schema.ServiceStatus) error {
	stmt, err := repo.db.Prepare(pgInsertServiceStatusQuery)
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

const pgSelectServiceStatusQuery = `
  SELECT
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at
  FROM dockmon_liveness_target WHERE service_name = $1`

// GetServiceStatus gets a specified service status from the database.
func (repo *PgServiceRepo) GetServiceStatus(serviceName string) (schema.ServiceStatus, error) {
	var s schema.ServiceStatus
	err := repo.db.QueryRow(pgSelectServiceStatusQuery, serviceName).Scan(
		&s.ServiceName, &s.LivenessURL, &s.LivenessInterval, &s.ShouldRestart,
		&s.FailAfter, &s.IsHealty, &s.Restarts, &s.ConsecutiveFailedHealthChecks,
		&s.LastRestarted, &s.LastHealthSuccess, &s.LastHealthFailure, &s.CreatedAt)
	if err != nil {
		return emptyServiceStatus, err
	}
	return s, err
}

const pgSelectServiceStatusesQuery = `
  SELECT
    service_name, liveness_url, liveness_interval, should_restart, fail_after,
    is_healty, number_of_restarts, consecutive_failed_health_checks,
    last_restarted, last_health_success, last_health_failure, created_at
  FROM dockmon_liveness_target ORDER BY service_name`

// GetServiceStatuses gets all service statuses from the database.
func (repo *PgServiceRepo) GetServiceStatuses() ([]schema.ServiceStatus, error) {
	rows, err := repo.db.Query(pgSelectServiceStatusesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return createServiceStatusesFromRows(rows)
}

// createServiceStatusesFromRows turns a resulting list of rows into
// a list of service statuses.
func createServiceStatusesFromRows(rows *sql.Rows) ([]schema.ServiceStatus, error) {
	statuses := make([]schema.ServiceStatus, 0)
	var s schema.ServiceStatus
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

const pgSetServiceSuccessQuery = `
  UPDATE dockmon_liveness_target SET
    last_health_success = $1, is_healty = TRUE, consecutive_failed_health_checks = 0
    WHERE service_name = $2`

// SaveHealthSuccess records a health check success for a given service.
func (repo *PgServiceRepo) SaveHealthSuccess(serviceName string, timestamp time.Time) error {
	stmt, err := repo.db.Prepare(pgSetServiceSuccessQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(timestamp, serviceName)
	return err
}

const pgSetServiceFailureQuery = `
  UPDATE dockmon_liveness_target SET
    last_health_failure = $1, is_healty = FALSE,
    consecutive_failed_health_checks = consecutive_failed_health_checks + 1
    WHERE service_name = $2`

// SaveHealthFailure records a health check failure for a given service.
func (repo *PgServiceRepo) SaveHealthFailure(serviceName string, timestamp time.Time) error {
	stmt, err := repo.db.Prepare(pgSetServiceFailureQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(timestamp, serviceName)
	return err
}

const pgSaveServiceRestartQuery = `
  UPDATE dockmon_liveness_target SET
    last_restarted = $1, consecutive_failed_health_checks = 0,
    number_of_restarts = number_of_restarts + 1
    WHERE service_name = $2`

// SaveRestart records a restart for a given service.
func (repo *PgServiceRepo) SaveRestart(serviceName string, timestamp time.Time) error {
	stmt, err := repo.db.Prepare(pgSaveServiceRestartQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(timestamp, serviceName)
	return err
}

// Close closes the underlying database connection.
func (repo *PgServiceRepo) Close() error {
	return repo.db.Close()
}
