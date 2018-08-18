export DOCKMON_DB_NAME=dockmon
export DOCKMON_DB_USER=dockmon
export DOCKMON_DB_HOST=localhost
export DOCKMON_DB_PASSWORD=password
export DOCKMON_DB_PORT=5432
export DOCKMON_USERNAME=this.guy
export DOCKMON_PASSWORD=password

build:
	go build

run-pg: build
	./dockmon -storage postgres

run-mysql: build
	DOCKMON_DB_PORT=3306 ./dockmon -storage mysql

run-sqlite: build
	DOCKMON_DB_NAME=dockmon.db ./dockmon -storage sqlite3

run-memory: build
	./dockmon -storage memory

clean:
	rm Gopkg.lock; dep ensure
