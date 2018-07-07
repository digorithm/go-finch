# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

RUN apt-get -y update
RUN apt-get -y install python3
RUN apt-get -y install python3-pip

RUN apt-get install -y \  
    libpng-dev \
    freetype* \
    libblas-dev \
    liblapack-dev \
    libatlas-base-dev \
    gfortran

ADD finchgo/requirements.txt /go/src/github.com/digorithm/meal_planner/finchgo/requirements.txt

RUN pip3 install -r src/github.com/digorithm/meal_planner/finchgo/requirements.txt

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
