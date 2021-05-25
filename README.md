# Disclaimer

This is a highly experimental work. Shouldn't be used anywhere near a production environment. The thesis describing in detail what's going on here is hosted here: https://open.library.ubc.ca/cIRcle/collections/ubctheses/24/items/1.0379346

The core of this theoretical/experimental work lives in the `finch/` directory. The other directories in this repo is virtually a simple CRUD application that makes use of the ideas in `finch/`.

**The code in this repo, and specially in `finch/` lives up to the meme about academics writing code that's hard to maintain.** Goes without saying that it doesn't reflect how I would build things that should be production-ready. Time was short, pressure was high, pay was low, and papers needed to be published at the time of the creation of `finch/` and this hasn't been updated in years. That said, proceed at your own caution. 


# Research and setup

Due to advancements in distributed systems and the increasing industrial demands placed on these systems, distributed systems are comprised of multiple complex
components (e.g databases and their replication infrastructure, caching components, proxies, and load balancers) each of which have their own complex configuration parameters that enable them to be tuned for given runtime requirements.

Software Engineers must manually tinker with many of these configuration parameters that change the behaviour and/or structure of the system in order to achieve their system requirements. In many cases, static configuration settings might not meet certain demands in a given context and ad hoc modifications of these configuration parameters can trigger unexpected behaviours, which can have negative effects on the quality of the overall system.

In this work, I show the design and analysis of Finch; a tool that injects a machine learning based MAPE-K feedback loop to existing systems to automate how these configuration parameters are set. Finch configures and optimizes the system to meet service-level agreements in uncertain workloads and usage patterns.
Rather than changing the core infrastructure of a system to fit the feedback loop, Finch asks the user to perform a small set of actions: instrumenting the code and configuration parameters, defining service-level objectives and agreements, and enabling programmatic changes to these configurations. As a result, Finch learns how to dynamically configure the system at runtime to self-adapt to its dynamic workloads.

I show how Finch can replace the trial-and-error engineering effort that otherwise would be spent manually optimizing a system’s wide array of configuration
parameters with an automated self-adaptive system.

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
