# Research and setup

## Running app and pg on docker compose

- Run `docker-compose up --build`

**Important:** If running for the first time, make sure init.sql has been executed on the container. Configure grafana /datasources pointing a datasource to prometheus:9090 and import the dash.json to have a complete dashboard

**Important**: If something goes wrong with the mapping of host to pgdata, try removing the db image with `docker rm db`.

## Monitoring

Instrumentation is being done with Prometheus + Granafa. After running docker compose up you can go to `localhost:3000` and login with admin:pass, you might have to setup the data source to Prometheus by passing the URI `http://prometheus:9090`.

Just in case here's a very good tutorial on setting all this up: https://finestructure.co/blog/2016/5/16/monitoring-with-prometheus-grafana-docker-part-1

## Big list of interesting metrics to monitor

Remember that for latency and other similar metrics, mean means nothing, focus on 99th percentile.

- System level:
  - CPU:
    - User
    - IOwait (A high iowait means that you are disk or network bound, This query produces a familiar value which is iowait as a percentage of CPU time averaged across all cores in the system and once graphed provides an insight into IO performance.): `avg(irate(node_cpu{job="node-exporter",mode="iowait"}[1m])) * 100`
    - Idle 
    - Note: Pass these functions to irate() to get the instant rate to see more details
  - Memory:
    - Mem useage in percentage: `((node_memory_MemTotal) - ((node_memory_MemFree+node_memory_Buffers+node_memory_Cached))) / node_memory_MemTotal * 100`
  - Disk
    - Still not sure about this

---

## Installation

1. Install PostgreSQL 9.4.x

2. Install Go 1.4.x, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

3. Create PostgreSQL database.
    ```
    cd $GOPATH/src/github.com/digorithm/meal_planner
    go get github.com/rnubel/pgmgr
    pgmgr db create
    ```

4. Run the PostgreSQL migration.
    ```
    pgmgr db migrate
    ```

5. Run the server
    ```
    cd $GOPATH/src/github.com/digorithm/meal_planner
    go run main.go
    ```


## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

* **DSN:** RDBMS database path. Default: `postgres://$(whoami)@localhost:5432/meal_planner?sslmode=disable`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.


## Running Migrations

Migration is handled by a separate project: [github.com/rnubel/pgmgr](https://github.com/rnubel/pgmgr).

Here's a quick tutorial on how to use it. For more details, read the tutorial [here](https://github.com/rnubel/pgmgr#usage).
```
# Installing the library
go get github.com/rnubel/pgmgr

# Create a new migration file
pgmgr migration {filename}

# Migrate all the way up
pgmgr db migrate

# Reset to the latest dump
pgmgr db drop
pgmgr db create
pgmgr db load

# Roll back the most recently applied migration, then run it again.
pgmgr db rollback
pgmgr db migrate

# Show the latest migration version
pgmgr db version
```