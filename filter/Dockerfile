ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="alpine:3.17.0"
FROM ${KINDEST_REPOSITORY}thavlik/bjjv-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY compiler/pkg/api compiler/pkg/api
COPY slideshow/pkg/api slideshow/pkg/api
COPY filter filter
WORKDIR /go/src/github.com/thavlik/bjjvb/filter/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o filter
FROM ${RUNTIME_BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/bjjvb/filter/cmd/filter /usr/local/bin
CMD ["filter"]
