FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/digorithm/meal_planner

ENV USER rodrigo
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET bginXRnaDjqwiOwb

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://rodrigo@localhost:5432/meal_planner?sslmode=disable&password=123

WORKDIR /go/src/github.com/digorithm/meal_planner

RUN godep go build

EXPOSE 8888
CMD ./meal_planner
