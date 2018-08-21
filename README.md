![dockmon-full-logo](https://user-images.githubusercontent.com/9406331/44312536-3455f480-a3fa-11e8-947d-a62a18e50d66.png)

# dockmon #
Liveness probing and and service restarting agent for docker containers on a single node. Allows users to specify any number of containers for which to conduct liveness probes, the frequency of health checks and optionally if a service should be restarted when it becomes unhealthy.

## Installation #
Pull the latest version of the docker image and run the service as follows:

```
docker run -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v serviceConf.yml:/etc/dockmon/serviceConf.yml
    -e DOCKMON_USERNAME=foo \
    -e DOCKMON_PASSWORD=bar \
    czarsimon/dockmon:1.0 -storage memory
```
This would start the dockmon service using the provided serviceConf.yml as the specification of what services to monitor and how. (See the usage section for information on how to write this specification).

Note the flag _**-storage memory**_ at the end of the command above. This tells dockmon to store the result of the health probes in memory. Other storage options are: postgres, mysql and sqlite3.  

Note: the `-v /var/run/docker.sock:/var/run/docker.sock` option can only be used in Linux environments.

## Usage #
In order for dockmon to have any value health check targets (refered to as services) has to be specified in a file name _serviceConf.yml_. Below is an example of what a serviceConf.yml file can look like:
```yaml
- serviceName: diplo-directory
  livenessUrl: http://localhost:1901/health
  livenessInterval: 10
  restart: true
  failAfter: 2
- serviceName: diplo-chat
  livenessUrl: http://localhost:1902/health
  livenessInterval: 15
  restart: true
  failAfter: 2
```

As seen above each service to monitor is specified with the following fields:
- _serviceName:_ Name of the service to monitor.
- _livenessUrl:_ URL to make the liveness probe to, the liveness probe will be a GET request which fill fail if the service returns a non 200 response.
- _livenessInterval:_ Time in seconds between liveness probes.
- _restart:_ Specifies if a service should be restarted if it is marked as unhealthy.
- _failAfter:_ Number of failed liveness probes required for the service to be marked as unhealthy.

Another option to providing the serviceConf.yml specification to dockmon by volume mounting `-v serviceConf.yml:/etc/dockmon/serviceConf.yml`, is to build your on docker image with serviceConf included. This can be done with a Dockerfile similar to this:
```Dockerfile
FROM czarsimon/dockmon:1.0
COPY serviceConf.yml /etc/dockmon/serviceConf.yml
```
### Storage options #
Dockmon has four options for storing the service health state as well as information such as number of restarts/liveness failures etc.

**Memory:** With this option dockmon's state is stored in memory and lost if it should be restarted.

**Sqlite3:** With this option an embeded sqlite database is set up and used for storing dockmon's state. The name of the database file is set by provideing the environment variable DOCKMON_DB_NAME. If a volume mapping is made to the sqlite database file, then the state of dockmon will survive restarts.

**Postgres:** Here a PostgreSQL database will be used to store the dockmon state. Connection information to an exteral Postgres database has to be specified by providing the environent variables: DOCKMON_DB_NAME, DOCKMON_DB_USER, DOCKMON_DB_HOST, DOCKMON_DB_PASSWORD and optionally DOCKMON_DB_PORT if not the default postgres port 5432 is used.

**Mysql:** Here a MySQL database will be used to store the dockmon state. As with the postgres option connection information has to be specified by providing the environent variables: DOCKMON_DB_NAME, DOCKMON_DB_USER, DOCKMON_DB_HOST, DOCKMON_DB_PASSWORD and optionally DOCKMON_DB_PORT if not the default mysql port 3306 is used.

Note: Database migrations will run when starting dockmon for the first time. Migration information will be stored in the table _dockmon_migrations_.

## Web UI #
Dockmon provides a web ui for inspecting the liveness status of the monitored services. By default the web ui can be accessed on port 7777 but this can be changed by setting the environment variable DOCKMON_PORT when starting dockmon.

Login credentials are required to access the web ui, these are the same that was set when starting dockmon by the environment variables DOCKMON_USERNAME and DOCKMON_PASSWORD.

### Screenshots #
| Login Page | Service List | Service Info |
|:-------------:|:-------:|:-------:|
|![dockmon-login](https://user-images.githubusercontent.com/9406331/44313173-c7475c80-a403-11e8-8087-7239b02f1709.png)|![dockmon-service-menu](https://user-images.githubusercontent.com/9406331/44313176-cd3d3d80-a403-11e8-8073-e9de25a3ae8e.png)|![dockmon-detailed-info](https://user-images.githubusercontent.com/9406331/44313171-c1ea1200-a403-11e8-80ff-97138d987f83.png)|

## CLI #
Another option to inspecting service status is to use the provided cli. Instal it by running: `go install github.com/CzarSimon/dockmon/cmd/cli/dockmon`

### Commands #
`$ dockmon get-services` lists all services monitored by dockmon along with a summary of their status.

`$ dockmon get-service [service-name]` displays the full status of a specified service.

`$ dockmon configure` prompts the user for configuration information such as remote host, username and password for the api.
