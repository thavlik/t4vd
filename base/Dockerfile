FROM golang:1.19.4
ENV GO111MODULE=on

WORKDIR /go/src/github.com/thavlik/bjjvb
COPY go.mod .
COPY go.sum .
RUN go mod download
