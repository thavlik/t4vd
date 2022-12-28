ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
ARG RUNTIME_BASE_IMAGE="debian:bullseye-slim"
FROM ${KINDEST_REPOSITORY}thavlik/bjjv-base:${KINDEST_TAG} AS builder
RUN apt-get update \
    && apt-get install -y libvpx-dev \
    && rm -rf /var/lib/apt/lists/*
COPY base/pkg base/pkg
COPY compiler/pkg/api compiler/pkg/api
COPY compiler/pkg/datastore/datastore.go compiler/pkg/datastore/datastore.go
COPY sources/pkg/api sources/pkg/api
COPY seer/pkg/api seer/pkg/api
COPY slideshow slideshow
WORKDIR /go/src/github.com/thavlik/bjjvb/slideshow/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o slideshow
FROM ${RUNTIME_BASE_IMAGE}
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
        libvpx-dev \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/github.com/thavlik/bjjvb/slideshow/cmd/slideshow /usr/local/bin
CMD ["slideshow"]
