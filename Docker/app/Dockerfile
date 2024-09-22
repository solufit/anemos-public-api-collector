FROM golang:1.23.0 AS base

RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . /go/src/app

RUN CGO_ENABLED=0 go build -o api main.go
RUN chmod +x api

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=base /go/src/app/api /app/api
ENTRYPOINT [ "/app/api" ]

