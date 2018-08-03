-- +migrate Up
CREATE TABLE dockmon_liveness_target (
  service_name VARCHAR(250) PRIMARY KEY,
  liveness_url VARCHAR(250),
  liveness_interval INTEGER,
  should_restart BOOLEAN,
  fail_after INTEGER,
  is_healty BOOLEAN,
  number_of_restarts INTEGER,
  consecutive_failed_health_checks INTEGER,
  last_restarted TIMESTAMP,
  last_health_success TIMESTAMP,
  last_health_failure TIMESTAMP,
  created_at TIMESTAMP
);
