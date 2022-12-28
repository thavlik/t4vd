ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/bjjv-base:${KINDEST_TAG} AS builder
COPY base/pkg base/pkg
COPY seer/pkg/api seer/pkg/api
COPY sources/pkg/api sources/pkg/api
COPY slideshow/pkg/api slideshow/pkg/api
COPY compiler compiler
WORKDIR /go/src/github.com/thavlik/bjjvb/compiler/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o compiler
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/github.com/thavlik/bjjvb/compiler/cmd/compiler /usr/local/bin
CMD ["compiler"]
