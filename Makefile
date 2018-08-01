export DOCKMON_DB_NAME=dockmon
export DOCKMON_DB_USER=dockmon
export DOCKMON_DB_HOST=localhost
export DOCKMON_DB_PASSWORD=password
export DOCKMON_DB_PORT=5432
export DOCKMON_USERNAME=this.guy
export DOCKMON_PASSWORD=password

build:
	go build

run-dev: build
	./dockmon --storage postgres
