package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
	"github.com/gobuffalo/packr"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
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
	case "sqlite3":
		return &SqliteServiceRepo{db: db}
	default:
		log.Fatalf("No ServiceRepository matching driver: %s\n", dbDriver)
		return nil
	}
}

func MigrateDB(dbDriver string, db *sql.DB) error {
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
