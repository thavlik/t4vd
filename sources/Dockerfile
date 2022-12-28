ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="alpine:3.17.0"
FROM ${KINDEST_REPOSITORY}thavlik/bjjv-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY base/cmd/iam base/cmd/iam
COPY compiler/pkg/api compiler/pkg/api
COPY seer/pkg/api seer/pkg/api
COPY sources sources
WORKDIR /go/src/github.com/thavlik/bjjvb/sources/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sources
FROM ${RUNTIME_BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/bjjvb/sources/cmd/sources /usr/local/bin
CMD ["sources"]
