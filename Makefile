export DOCKMON_DB_NAME=dockmon
export DOCKMON_DB_USER=dockmon
export DOCKMON_DB_HOST=localhost
export DOCKMON_DB_PASSWORD=password
export DOCKMON_DB_PORT=5432

build:
	go build

run-dev: build
	./dockmon --storage postgres
