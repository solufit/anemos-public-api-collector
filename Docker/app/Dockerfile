FROM golang:1.23.0

RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . /go/src/app

ENTRYPOINT ["go", "run", "main.go"]