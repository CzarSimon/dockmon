package datastore

import (
	"database/sql"
	"log"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
	_ "github.com/lib/pq"
)

var emptyServiceStatus schema.ServiceStatus = schema.ServiceStatus{}

// ServiceRepository interface to persist, update and retrieve service statuses.
type ServiceRepository interface {
	SaveService(serviceStatus schema.ServiceStatus) error
	GetServiceStatus(serviceName string) (schema.ServiceStatus, error)
	GetServiceStatuses() ([]schema.ServiceStatus, error)

	SaveHealthSuccess(serviceName string, timestamp time.Time) error
	SaveHealthFailure(serviceName string, timestamp time.Time) error
	SaveRestart(serviceName string, timestamp time.Time) error
	Close() error
}

func GetServiceRepository(dbDriver string, db *sql.DB) ServiceRepository {
	switch dbDriver {
	case "postgres":
		return &PgServiceRepo{db: db}
	case "mysql":
		return &MySQLServiceRepo{db: db}
	case "sqlite3":
		return &SqliteServiceRepo{db: db}
	default:
		log.Fatalf("No ServiceRepository matching driver: %s\n", dbDriver)
		return nil
	}
}
