-- +migrate Up
CREATE TABLE dockmon_liveness_target (
  service_name VARCHAR(150) PRIMARY KEY,
  liveness_url VARCHAR(250),
  liveness_interval INT,
  should_restart BOOLEAN,
  fail_after INT,
  is_healty BOOLEAN,
  number_of_restarts INT,
  consecutive_failed_health_checks INT,
  last_restarted TIMESTAMP,
  last_health_success TIMESTAMP,
  last_health_failure TIMESTAMP,
  created_at TIMESTAMP
);
