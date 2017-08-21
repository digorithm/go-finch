# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/digorithm/meal_planner

RUN go install github.com/digorithm/meal_planner

# Run the golang-docker command when the container starts.
# ENTRYPOINT /go/bin/meal_planner

ADD wait-for-postgres.sh /usr/local/bin/wait-for-postgres.sh

RUN ["chmod", "+x", "/usr/local/bin/wait-for-postgres.sh"]

CMD ["/usr/local/bin/wait-for-postgres.sh", "db:5432", "--", "/go/bin/meal_planner"]

# http server listens on port 8888.
EXPOSE 8888
