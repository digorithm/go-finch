# Research and setup

## Running app and pg on docker compose

- Run `docker-compose up --build`

**Important**: You need to manually run the script to build the database by ssh'ing into the pg container and pasting the code. I might automate it in the future, or not.

**Important 2:** if you want to run locally without Docker, go to `main.go` and change `DB_ADDR := "db"` to `DB_ADDR := "locahost"`

## Monitoring

Instrumentation is being done with Prometheus + Granafa. After running docker compose up you can go to `localhost:3000` and login with admin:pass, you might have to setup the data source to Prometheus by passing the URI `http://prometheus:9090`.

Just in case here's a very good tutorial on setting all this up: https://finestructure.co/blog/2016/5/16/monitoring-with-prometheus-grafana-docker-part-1

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