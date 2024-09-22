FROM golang:1.23.0

RUN go install github.com/amacneil/dbmate@latest

WORKDIR /go/src/app
COPY . /go/src/app
